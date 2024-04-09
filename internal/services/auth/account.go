package auth

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Onnywrite/grpc-auth/gen"
	"github.com/Onnywrite/grpc-auth/internal/lib/tokens"
	"github.com/Onnywrite/grpc-auth/internal/lib/ve"
	"github.com/Onnywrite/grpc-auth/internal/models"
	storage "github.com/Onnywrite/grpc-auth/internal/storage/common"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidIP          = errors.New("invalid IP")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrUserAlreadyRegistered = errors.New("user already exists")
	ErrUserDeleted           = errors.New("user has unregistred")

	ErrAlreadySignedUp = errors.New("you've already signed up and can log in")
	ErrSignedOut       = errors.New("you've signed out")
	ErrSignupNotExists = errors.New("signup does not exist")

	ErrAlreadyLoggedIn = errors.New("you've already logged in")

	ErrServiceNotExists = errors.New("service does not exist")

	ErrSessionNotExists         = errors.New("sesson does not exist")
	ErrSessionAlreadyTerminated = errors.New("session has already been terminated")

	ErrInternal     = errors.New("internal error")
	ErrUnauthorized = errors.New("unauthorized")
	ErrTokenExpired = errors.New("token has expired")
)

// Throws:
//
//	ValidationErrorsList
//	ErrUserAlreadyRegistered
//	ErrUserDeleted
//	ErrInternal
func (a *AuthService) Register(ctx context.Context, user *models.User, info models.SessionInfo) (*gen.IdTokens, error) {
	const op = "auth.AuthService.Register"
	log := a.log.With(slog.String("op", op))

	id := user.Idendifier()
	user.Login = &id.Value
	log = log.With(slog.String("login_type", id.Key), slog.String("login", id.Value))

	if err := validator.New().Struct(user); err != nil {
		errs := ve.From(err.(validator.ValidationErrors))
		log.Error("validation error", slog.String("validation_errors", errs.Error()))
		return nil, errs
	}
	log.Info("passed validation")

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("cannot hash password", slog.String("error", err.Error()))
		return nil, ErrInternal
	}
	user.Password = string(hashed)

	saved, err := a.saveUser(ctx, user)
	if err != nil {
		log.Error("saving user error", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("user registred", slog.Int64("id", saved.Id))

	// TODO: generate tokens:

	// Create signup (service_id = 0)
	// Log in this signup
	// send token, recieved by Login

	return &gen.IdTokens{
		IdToken: "TODO",
		Profile: &gen.UserProfile{
			Id:    saved.Id,
			Login: saved.Login,
			Email: saved.Email,
			Phone: saved.Phone,
		},
	}, nil
}

// Throws:
//
//	ErrInvalidCredentials
//	ErrServiceNotExists
//	ErrSignedOut
//	ErrAlreadySignedUp
//	ErrInternal
func (a *AuthService) Signup(ctx context.Context, identifier models.UserIdentifier, serviceId int64) error {
	const op = "auth.AuthService.Signup"
	log := a.log.With(slog.String("op", op), slog.Any("identifier", identifier), slog.Int64("service_id", serviceId))

	// TODO: getUser wrapper
	user, err := a.db.User(ctx, identifier)
	if errors.Is(err, storage.ErrEmptyResult) {
		log.Error("invalid credentials", slog.String("error", err.Error()))
		return ErrInvalidCredentials
	}
	if err != nil {
		log.Error("internal error", slog.String("error", err.Error()))
		return ErrInternal
	}
	log = log.With(slog.Int64("user_id", user.Id))
	log.Info("user found")

	if user.IsDeleted() {
		log.Error("user deleted")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(identifier.Password))
	if err != nil {
		log.Error("invalid password", slog.String("error", err.Error()))
		return ErrInvalidCredentials
	}
	log.Info("password and hash match")

	signup := models.Signup{
		UserId:    user.Id,
		ServiceId: serviceId,
	}
	su, err := a.db.SaveSignup(ctx, signup)
	switch {
	case errors.Is(err, storage.ErrUniqueConstraint):
		log.Info("checking if user's signed out")
		su, _ = a.db.Signup(ctx, signup.UserId, signup.ServiceId)
		if su.IsDeleted() {
			log.Error("signup deleted")
			return ErrSignedOut
		}
		log.Error("signup exists")
		return ErrAlreadySignedUp
	case errors.Is(err, storage.ErrFKConstraint):
		log.Error("service with this id not found")
		return ErrServiceNotExists
	case err != nil:
		log.Error("failed to save signup", slog.String("error", err.Error()))
		return ErrInternal
	}
	log.Info("signed up", slog.Int64("signup_id", su.Id))

	return nil
}

// Throws:
//
//	ValidationErrorsList
//	ErrInvalidCredentials
//	ErrSignupNotExists
//	ErrAlreadyLoggedIn
//	ErrInternal
//
// TODO: need rollback (delete session, if an error occured while creating tokens)
func (a *AuthService) Login(ctx context.Context, identifier models.UserIdentifier, sessionInfo models.SessionInfo, serviceId int64) (*models.Tokens, error) {
	const op = "auth.AuthService.Login"
	log := a.log.With(slog.String("op", op), slog.Any("identifier", identifier), slog.Any("session", sessionInfo))

	if err := validator.New().Struct(sessionInfo); err != nil {
		errs := ve.From(err.(validator.ValidationErrors))
		log.Error("validation error", slog.String("validation_errors", errs.Error()))
		return nil, errs
	}

	user, err := a.db.User(ctx, identifier)
	if errors.Is(err, storage.ErrEmptyResult) {
		log.Error("invalid credentials", slog.String("error", err.Error()))
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		log.Error("internal error", slog.String("error", err.Error()))
		return nil, ErrInternal
	}
	log = log.With(slog.Int64("user_id", user.Id))
	log.Info("user found")

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(identifier.Password))
	if err != nil {
		log.Error("invalid password", slog.String("error", err.Error()))
		return nil, ErrInvalidCredentials
	}
	log.Info("password and hash match")

	su, err := a.db.Signup(ctx, user.Id, serviceId)
	if err != nil || su.IsDeleted() {
		log.Error("signup not found")
		return nil, ErrSignupNotExists
	}
	log = log.With(slog.Int64("signup_id", su.Id))
	log.Info("user signed up")

	session := &models.Session{
		UserId:    su.UserId,
		ServiceId: su.ServiceId,
		Info:      sessionInfo,
	}
	saved, err := a.db.SaveSession(ctx, session)
	if errors.Is(err, storage.ErrUniqueConstraint) {
		log.Error("could not save session", slog.String("error", err.Error()))
		// TODO: not done
		a.checkIfSessionTerminated(ctx, session)
		return nil, ErrAlreadyLoggedIn
	}
	if err != nil {
		log.Error("internal error", slog.String("error", err.Error()))
		return nil, ErrInternal
	}
	log = log.With(slog.String("session_uuid", saved.UUID))
	log.Info("saved session")

	refresh, err := tokens.Refresh(&models.RefreshToken{
		SessionUUID: saved.UUID,
		Exp:         time.Now().Add(a.tokenTTL).Unix(),
	})
	if err != nil {
		log.Error("cannot create refresh token", slog.String("error", err.Error()))
		return nil, ErrInternal
	}

	access, err := tokens.Access(&models.AccessToken{
		Id:        user.Id,
		Login:     user.Login,
		ServiceId: serviceId,
		Roles:     []string{},
		Exp:       time.Now().Add(a.refreshTokenTTL).Unix(),
	})
	if err != nil {
		log.Error("cannot create refresh token", slog.String("error", err.Error()))
		return nil, ErrInternal
	}

	log.Info("logged in successfully")

	return &models.Tokens{
		Refresh: refresh,
		Access:  access,
	}, nil
}

func (a *AuthService) Logout(ctx context.Context, refresh string) error {
	const op = "auth.AuthService.Logout"
	log := a.log.With(slog.String("op", op))

	token, err := tokens.ParseRefresh(refresh)
	if errors.Is(err, tokens.ErrTokenExpired) {
		return ErrTokenExpired
	}
	if err != nil {
		log.Error("could not process refresh token", slog.String("token", refresh), slog.String("error", err.Error()))
		return ErrUnauthorized
	}
	log = log.With(slog.String("session_uuid", token.SessionUUID))
	log.Info("token is processed")

	err = a.db.TerminateSession(ctx, token.SessionUUID)
	if errors.Is(err, storage.ErrEmptyResult) {
		log.Error("checking if session terminated", slog.String("error", err.Error()))
		return a.checkIfSessionTerminatedById(ctx, token.SessionUUID)
	}
	if err != nil {
		log.Error("could not terminate session", slog.String("error", err.Error()))
		return err
	}
	log.Info("logged out successfully")

	return nil
}

// Throws;
//
//	ErrInvalidCredentials if 'login' or password is invalid
//	ErrUserDeleted if user has deleted their account
//	ErrInternal in any unexpected situation
func (a *AuthService) getUser(ctx context.Context, identifier models.UserIdentifier) (*models.SavedUser, error) {
	const op = "auth.AuthService.getUser"
	log := a.log.With(slog.String("op", op), slog.String("login_type", identifier.Key), slog.String("login", identifier.Value))

	user, err := a.db.User(ctx, identifier)
	if errors.Is(err, storage.ErrEmptyResult) {
		log.Error("invalid identifier", slog.String("error", err.Error()))
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		log.Error("internal error", slog.String("error", err.Error()))
		return nil, ErrInternal
	}
	log = log.With(slog.Int64("user_id", user.Id))
	log.Info("user found")

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(identifier.Password))
	if err != nil {
		log.Error("invalid password", slog.String("error", err.Error()))
		return nil, ErrInvalidCredentials
	}
	log.Info("password and hash match")

	if user.IsDeleted() {
		log.Error("user deleted", slog.Time("deleted_at", *user.DeletedAt))
		return nil, ErrUserDeleted
	}

	return user, nil
}

func (a *AuthService) saveUser(ctx context.Context, user *models.User) (*models.SavedUser, error) {
	const op = "auth.AuthService.saveUser"
	log := a.log.With(slog.String("op", op))

	u, err := a.db.SaveUser(ctx, user)
	if errors.Is(err, storage.ErrUniqueConstraint) {
		u, err = a.getUser(ctx, *user.Idendifier())
		if errors.Is(err, ErrUserDeleted) {
			log.Error("user deleted", slog.Int64("id", u.Id))
			return nil, err
		}
		return nil, ErrUserAlreadyRegistered
	}
	if err != nil {
		log.Error("saving error", slog.String("error", err.Error()))
		return nil, ErrInternal
	}

	return u, nil
}
