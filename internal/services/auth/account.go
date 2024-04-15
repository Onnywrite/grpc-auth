package auth

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/Onnywrite/grpc-auth/gen"
	se "github.com/Onnywrite/grpc-auth/internal/lib/service-errors"
	"github.com/Onnywrite/grpc-auth/internal/lib/tokens"
	"github.com/Onnywrite/grpc-auth/internal/models"
	auth "github.com/Onnywrite/grpc-auth/internal/services/auth/common"
	storage "github.com/Onnywrite/grpc-auth/internal/storage/common"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func identify(u *models.User) {
	switch {
	case u.Login != nil:
		return
	case u.Email != nil:
		u.Login = u.Email
		return
	case u.Phone != nil:
		u.Login = u.Phone
		return
	}
}

// Throws:
//
//	Errors
//	ErrUserAlreadyRegistered
//	ErrUserDeleted
//	ErrInternal
func (a *AuthService) Register(ctx context.Context, user *models.User, info models.SessionInfo) (*gen.IdTokens, error) {
	const op = "auth.AuthService.Register"
	log := a.log.With(slog.String("op", op))

	identify(user)

	if err := validator.New().Struct(user); err != nil {
		errs := se.From(err.(validator.ValidationErrors))
		log.Error("user validation error", slog.String("validation_errors", errs.Error()))
		return nil, errs
	}
	log.Info("passed validation")

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("cannot hash password", slog.String("error", err.Error()))
		return nil, auth.ErrInternal
	}
	user.Password = string(hashed)

	saved, err := a.db.SaveUser(ctx, user)
	if err != nil {
		log.Error("saving user error", slog.String("error", err.Error()))
		return nil, err
	}
	log.Info("user registred", slog.Int64("id", saved.Id))

	return a.openIdConnect(saved)
}

// Throws:
//
//	Errors
//	ErrInvalidCredentials
//	ErrSignupNotExists
//	ErrAlreadyLoggedIn
//	ErrInternal
func (a *AuthService) Login(ctx context.Context, user *models.User, info models.SessionInfo) (*gen.IdTokens, error) {
	const op = "auth.AuthService.Login"
	log := a.log.With(slog.String("op", op), slog.Any("session", info))

	if err := validator.New().StructExcept(user, "Login"); err != nil {
		errs := se.From(err.(validator.ValidationErrors))
		log.Error("user validation error", slog.String("validation_errors", errs.Error()))
		return nil, errs
	}
	log.Info("user passed validation")

	if err := validator.New().Struct(info); err != nil {
		errs := se.From(err.(validator.ValidationErrors))
		log.Error("session info validation error", slog.String("validation_errors", errs.Error()))
		return nil, errs
	}
	log.Info("session info passed validation")

	u := &models.SavedUser{}
	var err error
	switch {
	case user.Login != nil:
		u, err = a.db.UserByLogin(ctx, *user.Login)
	case user.Email != nil:
		u, err = a.db.UserByEmail(ctx, *user.Email)
	case user.Phone != nil:
		u, err = a.db.UserByPhone(ctx, *user.Phone)
	}
	if err != nil {
		return nil, err
	}
	log = log.With(slog.Int64("user_id", u.Id))
	log.Info("user found")

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		log.Error("invalid password", slog.String("error", err.Error()))
		return nil, auth.ErrInvalidCredentials
	}
	log.Info("password and hash match")

	return a.openIdConnect(u)
}

func (a *AuthService) openIdConnect(user *models.SavedUser) (*gen.IdTokens, error) {
	const op = "a.AuthService.openIdConnect"
	log := a.log.With(slog.String("op", op))

	tkn, err := tokens.Id(models.IdToken{
		Id:  user.Id,
		Iss: "sso.onnywrite.com",
		Sub: "sso.onnywrite.com",
		Exp: time.Now().Add(a.idTokenTTL).Unix(),
	})
	if err != nil {
		log.Error("could not generate id token", slog.String("error", err.Error()))
		return nil, auth.ErrInternal
	}
	log.Info("opened id connect successfully")

	return &gen.IdTokens{
		IdToken: tkn,
		Profile: &gen.UserProfile{
			Id:    user.Id,
			Login: user.Login,
			Email: user.Email,
			Phone: user.Phone,
		},
	}, nil
}

// Throws:
//
//	Errors
//	ErrInvalidCredentials
//	ErrServiceNotExists
//	ErrSignedOut
//	ErrAlreadySignedUp
//	ErrInternal
func (a *AuthService) Signup(ctx context.Context, idToken string, serviceId int64, info models.SessionInfo) (*gen.AppTokens, error) {
	const op = "auth.AuthService.Signup"
	log := a.log.With(slog.String("op", op), slog.Int64("service_id", serviceId), slog.Any("session", info))

	if err := validator.New().Struct(info); err != nil {
		errs := se.From(err.(validator.ValidationErrors))
		log.Error("session info validation error", slog.String("validation_errors", errs.Error()))
		return nil, errs
	}

	token, err := tokens.ParseId(idToken)
	if errors.Is(err, tokens.ErrTokenExpired) {
		return nil, auth.ErrTokenExpired
	}
	if err != nil {
		log.Error("could not parse id token", slog.String("error", err.Error()))
		return nil, auth.ErrInternal
	}

	user, err := a.db.UserById(ctx, token.Id)
	if err != nil {
		return nil, err
	}

	su, err := a.db.SaveSignup(ctx, models.Signup{
		UserId:    user.Id,
		ServiceId: serviceId,
	})
	if err != nil {
		return nil, err
	}

	log.Info("signed up")

	return a.openSession(ctx, su, user, info)
}

func (a *AuthService) Signin(ctx context.Context, idToken string, serviceId int64, info models.SessionInfo) (*gen.AppTokens, error) {
	const op = "auth.AuthService.Signin"
	log := a.log.With(slog.String("op", op), slog.Int64("service_id", serviceId), slog.Any("session", info))

	if err := validator.New().Struct(info); err != nil {
		errs := se.From(err.(validator.ValidationErrors))
		log.Error("session info validation error", slog.String("validation_errors", errs.Error()))
		return nil, errs
	}
	log.Info("passed validation")

	token, err := tokens.ParseId(idToken)
	if errors.Is(err, tokens.ErrTokenExpired) {
		return nil, auth.ErrTokenExpired
	}
	if err != nil {
		log.Error("could not parse id token", slog.String("error", err.Error()))
		return nil, auth.ErrInternal
	}

	user, err := a.db.UserById(ctx, token.Id)
	if err != nil {
		return nil, err
	}

	su, err := a.db.SignupByServiceAndUser(ctx, serviceId, token.Id)
	if err != nil {
		return nil, err
	}

	log.Info("openning session")
	return a.openSession(ctx, su, user, info)
}

func (a *AuthService) openSession(ctx context.Context, signup *models.SavedSignup, user *models.SavedUser, info models.SessionInfo) (*gen.AppTokens, error) {
	const op = "a.AuthService.openSession"
	log := a.log.With(slog.String("op", op), slog.Int64("service_id", signup.ServiceId), slog.Int64("user_id", signup.UserId))

	session, rollback, err := a.db.SaveSession(ctx, &models.Session{
		ServiceId: signup.ServiceId,
		UserId:    signup.UserId,
		Info:      info,
	})
	if err != nil {
		return nil, err
	}

	var refresh, access string
	errorsCh := make(chan error, 5)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		var err33 error
		access, err33 = a.updateAccessToken(ctx, user, signup.ServiceId)
		errorsCh <- err33
	}()

	go func() {
		defer wg.Done()
		var err33 error
		refresh, err33 = tokens.Refresh(&models.RefreshToken{
			SessionUUID: session.UUID,
			Exp:         time.Now().Add(a.refreshTokenTTL).Unix(),
		})
		if err33 != nil {
			log.Error("could not generate refresh token", slog.String("error", err.Error()))
			errorsCh <- auth.ErrInternal
			return
		}
		errorsCh <- nil
	}()

	wg.Wait()
	for err := range errorsCh {
		if err != nil {
			if err2 := rollback(); err2 != nil {
				log.Error("cannot rollback session saving", slog.String("error", err.Error()))
			}
			return nil, err
		}
	}

	return &gen.AppTokens{
		Access:  access,
		Refresh: refresh,
		Profile: &gen.UserProfile{
			Id:    user.Id,
			Login: user.Login,
			Email: user.Email,
			Phone: user.Phone,
		},
	}, nil
}

func (a *AuthService) updateAccessToken(ctx context.Context, user *models.SavedUser, serviceId int64) (string, error) {
	const op = "a.AuthService.updateAccessToken"
	log := a.log.With(slog.String("op", op))

	// TODO: get roles

	access, err := tokens.Access(&models.AccessToken{
		Id:        user.Id,
		Login:     user.Login,
		ServiceId: serviceId,
		Roles:     []string{},
		Exp:       time.Now().Add(a.tokenTTL).Unix(),
	})
	if err != nil {
		log.Error("could not generate access token", slog.String("error", err.Error()))
		return "", auth.ErrInternal
	}

	return access, nil
}

func (a *AuthService) Signout(ctx context.Context, refresh string) error {
	const op = "auth.AuthService.Signout"
	log := a.log.With(slog.String("op", op))

	token, err := tokens.ParseRefresh(refresh)
	if errors.Is(err, tokens.ErrTokenExpired) {
		return auth.ErrTokenExpired
	}
	if err != nil {
		log.Error("could not process refresh token", slog.String("token", refresh), slog.String("error", err.Error()))
		return auth.ErrUnauthorized
	}
	log = log.With(slog.String("session_uuid", token.SessionUUID))
	log.Info("token is processed")

	err = a.db.TerminateSession(ctx, token.SessionUUID)
	if errors.Is(err, storage.ErrEmptyResult) {
		log.Error("session does not exist or already terminated", slog.String("error", err.Error()))
		return auth.ErrSessionTerminated
	}
	if err != nil {
		log.Error("could not terminate session", slog.String("error", err.Error()))
		return err
	}
	log.Info("signed out successfully")

	return nil
}

func (a *AuthService) Resignin(ctx context.Context, refresh string) (*gen.AppTokens, error) {
	const op = "a.AuthService.Resignin"
	log := a.log.With(slog.String("op", op))

	token, err := tokens.ParseRefresh(refresh)
	if errors.Is(err, tokens.ErrTokenExpired) {
		return nil, auth.ErrTokenExpired
	}
	if err != nil {
		log.Error("could not process refresh token", slog.String("token", refresh), slog.String("error", err.Error()))
		return nil, auth.ErrUnauthorized
	}
	log = log.With(slog.String("session_uuid", token.SessionUUID))
	log.Info("token is processed")

	session, err := a.db.SessionByUuid(ctx, token.SessionUUID)
	if err != nil {
		return nil, err
	}
	log.Info("got session")

	user, err := a.db.UserById(ctx, session.UserId)
	if err != nil {
		return nil, err
	}
	log.Info("got user")

	access, err := a.updateAccessToken(ctx, user, session.ServiceId)
	if err != nil {
		return nil, err
	}
	log.Info("resigned in")

	return &gen.AppTokens{
		Access:  access,
		Refresh: refresh,
		Profile: &gen.UserProfile{
			Id:    user.Id,
			Login: user.Login,
			Email: user.Email,
			Phone: user.Phone,
		},
	}, nil
}

func (a *AuthService) Check(ctx context.Context, access string) error {
	const op = "a.AuthService.Check"
	log := a.log.With(slog.String("op", op))

	token, err := tokens.ParseAccess(access)
	if errors.Is(err, tokens.ErrTokenExpired) {
		return auth.ErrTokenExpired
	}
	if err != nil {
		log.Error("could not process refresh token", slog.String("token", access), slog.String("error", err.Error()))
		return auth.ErrUnauthorized
	}
	log.Info("token checked", slog.Int64("user_id", token.Id))

	return nil
}
