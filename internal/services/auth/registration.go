package auth

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/Onnywrite/grpc-auth/internal/storage"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenEnv = "TOKEN_SECRET"
)

var (
	ErrUserRegistred      = errors.New("user already exists")
	ErrUserNotExists      = errors.New("user does not exist")
	ErrPasswordTooLong    = errors.New("password is longer than 72 bytes")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAlreadySignedUp    = errors.New("you've already signed up and can log in")
	ErrNoSuchService      = errors.New("service does not exist")
	ErrInternal           = errors.New("internal error")
)

func (a *AuthServiceImpl) Register(ctx context.Context, optLogin, optEmail, optPhone *string, password string) error {
	const op = "auth.AuthServiceImpl.Register"
	log := a.log.With(slog.String("op", op))

	log.Debug("switching login")
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
	log.Info("switched login")

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

	saved, err := a.db.SaveUser(ctx, user)
	if errors.Is(err, storage.ErrUserExists) {
		return ErrUserRegistred
	}
	if err != nil {
		log.Error("saving error", slog.String("error", err.Error()))
	}
	log = log.With(slog.Int64("id", saved.Id))
	log.Info("user registred")

	return nil
}

func (a *AuthServiceImpl) SignUp(ctx context.Context, key, password string, serviceId int64) error {
	const op = "auth.AuthServiceImpl.SignUp"
	log := a.log.With(slog.String("op", op), slog.String("key", key), slog.Int64("service_id", serviceId))

	user, err := a.recognize(ctx, key)
	if err != nil {
		log.Error("recognition failed")
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Error("invalid password", slog.String("error", err.Error()))
		return ErrInvalidCredentials
	}

	err = a.db.SaveSignup(ctx, models.Signup{
		UserId:    user.Id,
		ServiceId: serviceId,
	})
	if errors.Is(err, storage.ErrNoSuchPrimaryKey) {
		log.Error("falied to save signup", slog.String("error", err.Error()))
		return ErrUserNotExists
	}
	if err != nil {
		log.Error("failed to save signup", slog.String("error", err.Error()))
		return ErrInternal
	}

	return nil
}

func (a *AuthServiceImpl) LogIn(ctx context.Context,
	key, password string,
	serviceId int64 /*some session data*/) (refresh string, access string, err error) {
	const op = "auth.AuthServiceImpl.LogIn"
	log := a.log.With(slog.String("op", op), slog.String("key", key), slog.Int64("service_id", serviceId))

	user, err := a.recognize(ctx, key)
	if err != nil {
		log.Error("recognition failed")
		return "", "", err
	}

	a.db.Signup

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         user.Id,
		"login":      user.Login,
		"service_id": serviceId,
		// probably custom type like models.Role
		"roles": []string{},
		"exp":   time.Now().Add(a.tokenTTL).Unix(),
	})

	token, err := tkn.SignedString(os.Getenv(tokenEnv))
	if err != nil {
		log.Error("failed to sign token", slog.String("error", err.Error()))
		return "", ErrInternal
	}

	return token, nil
}

func (a *AuthServiceImpl) LogOut(ctx context.Context, token string) error {
	return nil
}

func (a *AuthServiceImpl) recognize(ctx context.Context, key string) (user *models.SavedUser, err error) {
	const op = "auth.AuthServiceImpl.recognize"
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
	if errors.Is(err, storage.ErrUserNotFound) {
		log.Error("recieving user", slog.String("error", err.Error()))
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		log.Error("recieving user", slog.String("error", err.Error()))
		return nil, ErrInternal
	}
	log.Info("switched key and recognized user", slog.Int64("id", user.Id))

	return user, nil
}
