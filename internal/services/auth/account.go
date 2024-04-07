package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/netip"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/lib/tokens"
	"github.com/Onnywrite/grpc-auth/internal/models"
	storage "github.com/Onnywrite/grpc-auth/internal/storage/common"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidIP          = errors.New("invalid IP")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrUserAlreadyRegistered = errors.New("user already exists")
	ErrUserNotExists         = errors.New("user does not exist")

	ErrAlreadySignedUp = errors.New("you've already signed up and can log in")
	ErrSignupDeleted   = errors.New("you've signed out")
	ErrSignupNotExists = errors.New("signup does not exist")

	ErrAlreadyLoggedIn = errors.New("you've already logged in")

	ErrServiceNotExists = errors.New("service does not exist")

	ErrSessionNotExists         = errors.New("sesson does not exist")
	ErrSessionAlreadyTerminated = errors.New("session has already been terminated")

	ErrInternal     = errors.New("internal error")
	ErrUnauthorized = errors.New("unauthorized")
	ErrTokenExpired = errors.New("token has expired")
)

func (a *AuthService) Register(ctx context.Context, optLogin, optEmail, optPhone *string, password string) error {
	const op = "auth.AuthService.Register"
	log := a.log.With(slog.String("op", op))

	log.Debug("switching login type")
	var login string
	switch {
	case optLogin != nil:
		log = log.With(slog.String("login_type", "login"))
		log.Debug("choose login")
		login = *optLogin
	case optEmail != nil:
		log = log.With(slog.String("login_type", "email"))
		log.Debug("choose email")
		login = *optEmail
	case optPhone != nil:
		log = log.With(slog.String("login_type", "phone"))
		log.Debug("choose phone")
		login = *optPhone
	}
	log = log.With(slog.String("login", login))
	log.Info("switched login type")

	user := &models.User{
		Login:    login,
		Email:    optEmail,
		Phone:    optPhone,
		Password: password,
	}

	if err := validator.New().Struct(user); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			log.Error("validation error", slog.String("errors", errs.Error()))
			return ErrInvalidCredentials
		}
		log.Error("validation error", slog.String("error", err.Error()))
		return ErrInvalidCredentials
	}
	log.Info("passed validation")

	saved, err := a.db.SaveUser(ctx, user)
	if errors.Is(err, storage.ErrUniqueConstraint) {
		return ErrUserAlreadyRegistered
	}
	if err != nil {
		log.Error("saving error", slog.String("error", err.Error()))
	}
	log.Info("user registred", slog.Int64("id", saved.Id))

	return nil
}

func (a *AuthService) Signup(ctx context.Context, key, password string, serviceId int64) error {
	const op = "auth.AuthService.Signup"
	log := a.log.With(slog.String("op", op), slog.String("key", key), slog.Int64("service_id", serviceId))

	user, err := a.recognize(ctx, key)
	if err != nil {
		log.Error("recognition failed")
		return err
	}
	log = log.With(slog.Int64("user_id", user.Id))
	log.Info("user recognized")

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
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
			return ErrSignupDeleted
		}
		log.Error("signup exists")
		return ErrAlreadySignedUp
	case err != nil:
		log.Error("failed to save signup", slog.String("error", err.Error()))
		return ErrInternal
	}
	log.Info("signed up", slog.Int64("signup_id", su.Id))

	return nil
}

func (a *AuthService) Login(ctx context.Context, key, password string, serviceId int64, browser, os, ip string) (*models.Tokens, error) {
	const op = "auth.AuthService.Login"
	log := a.log.With(slog.String("op", op), slog.String("key", key),
		slog.Int64("service_id", serviceId), slog.String("ip", ip), slog.String("os", os), slog.String("browser", browser))

	user, err := a.recognize(ctx, key)
	if err != nil {
		log.Error("recognition failed")
		return nil, err
	}
	log = log.With(slog.Int64("user_id", user.Id))
	log.Info("user recognized")

	su, err := a.db.Signup(ctx, user.Id, serviceId)
	if err != nil {
		log.Error("signup not found")
		return nil, ErrSignupNotExists
	}
	log = log.With(slog.Int64("signup_id", su.Id))
	log.Info("user signed up")

	ipParsed, err := netip.ParseAddr(ip)
	if err != nil {
		log.Error("could not parse IP address")
		return nil, ErrInvalidIP
	}

	session, err := a.db.SaveSession(ctx, &models.Session{
		UserId:    su.UserId,
		ServiceId: su.ServiceId,
		IP:        ipParsed,
		Browser:   &browser,
		OS:        &os,
	})
	if err != nil {
		log.Error("could not save session", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrUniqueConstraint) {
			return nil, ErrAlreadyLoggedIn
		}
		return nil, ErrInternal
	}
	log = log.With(slog.String("session_uuid", session.UUID.String()))
	log.Info("saved session")

	refresh, err := tokens.Refresh(&models.RefreshToken{
		SessionUUID: session.UUID,
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
	log = log.With(slog.String("session_uuid", token.SessionUUID.String()))
	log.Info("token is processed")

	err = a.db.TerminateSession(ctx, token.SessionUUID)
	if errors.Is(err, storage.ErrEmptyResult) {
		log.Error("checking if session terminated", slog.String("error", err.Error()))
		return a.checkIfSessionTerminated(ctx, token.SessionUUID)
	}
	if err != nil {
		log.Error("could not terminate session", slog.String("error", err.Error()))
		return err
	}
	log.Info("logged out successfully")

	return nil
}

func (a *AuthService) checkIfSessionTerminated(ctx context.Context, uuid uuid.UUID) error {
	const op = "auth.AuthService.checkIfSessionTerminated"
	log := a.log.With("op", op)

	session, err := a.db.Session(ctx, uuid)
	if err != nil {
		log.Error("session has been deleted", slog.String("error", err.Error()))
		return ErrSessionNotExists
	}
	if session.IsTerminated() {
		log.Error("session already terminated", slog.String("error", err.Error()))
		return ErrSessionAlreadyTerminated
	}

	return nil
}

func (a *AuthService) recognize(ctx context.Context, key string) (user *models.SavedUser, err error) {
	const op = "auth.AuthService.recognize"
	log := a.log.With(slog.String("op", op))

	validate := validator.New()
	log.Debug("switching key")
	switch {
	case validate.Var(key, "email") == nil:
		log = log.With(slog.String("key_type", "email"))
		user, err = a.db.UserByEmail(ctx, key)
	case validate.Var(key, "email") == nil:
		log = log.With(slog.String("key_type", "phone"))
		user, err = a.db.UserByPhone(ctx, key)
	default:
		log = log.With(slog.String("key_type", "login"))
		user, err = a.db.UserByLogin(ctx, key)
	}
	if errors.Is(err, storage.ErrEmptyResult) {
		log.Error("invalid credentials", slog.String("error", err.Error()))
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		log.Error("internal error", slog.String("error", err.Error()))
		return nil, ErrInternal
	}
	log.Info("switched key and recognized user", slog.Int64("id", user.Id))

	return user, nil
}
