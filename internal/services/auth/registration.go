package auth

import (
	"context"
	"errors"
	"fmt"
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
	ErrPasswordTooLong    = errors.New("password is longer than 72 bytes")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAlreadySignedUp    = errors.New("you've already signed up and can log in")
	ErrNoSuchService      = errors.New("service does not exist")
)

func (a *AuthServiceImpl) Register(ctx context.Context, optLogin, optEmail, optPhone, password string) error {
	const op = "auth.AuthServiceImpl.Register"
	logger := slog.With(slog.String("op", op))

	optLogin, err := initLogin(optLogin, optEmail, optPhone)
	if err != nil {
		logger.Error("cannot init login",
			slog.String("error", err.Error()),
			slog.String("login", optLogin),
			slog.String("email", optEmail),
			slog.String("phone", optPhone),
		)
		return err
	}

	logger = slog.With(slog.String("login", optLogin))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("password is too long", slog.String("error", err.Error()), slog.String("password", password))
		return ErrPasswordTooLong
	}

	user := &models.User{
		Login:    &optLogin,
		Email:    nilIfEmpty(optEmail),
		Phone:    nilIfEmpty(optPhone),
		Password: string(passHash),
	}

	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		logger.Error("validation errors", slog.Any("errors", errs.Error()))
		return ErrInvalidCredentials
	}

	_, err = a.db.SaveUser(ctx, user)
	if err != nil {
		logger.Error("cannot save", slog.String("error", err.Error()))
		return err
	}

	logger.Info("user signed up")
	return err
}

func (a *AuthServiceImpl) SignUp(ctx context.Context, key, password string, serviceId int64) (token string, err error) {
	const op = "auth.AuthServiceImpl.LogIn"
	logger := slog.With(slog.String("op", op), slog.String("key", key))

	user, err := a.userByKey(ctx, logger, key)
	if err != nil {
		return "", err
	}

	err = a.db.SaveSignup(ctx, models.Signup{
		UserId:    user.Id,
		ServiceId: serviceId,
	})
	if err != nil {
		logger.Error("cannot signup",
			slog.String("error", err.Error()),
			slog.Int64("service_id", serviceId),
			slog.Int64("user_id", user.Id))

		if errors.Is(err, storage.ErrSignupExists) {
			return "", ErrAlreadySignedUp
		}
		if errors.Is(err, storage.ErrNoSuchPrimaryKey) {
			return "", ErrNoSuchService
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.Error("hash and password are different")
		return "", ErrInvalidCredentials
	}

	// TODO: SaveLogin here

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         user.Id,
		"login":      user.Login,
		"email":      user.Email,
		"phone":      user.Phone,
		"service_id": serviceId,
		"exp":        time.Now().Add(a.tokenTTL).Unix(),
	})

	token, err = jwtToken.SignedString([]byte(os.Getenv(tokenEnv)))
	if err != nil {
		logger.Error("cannot sign JWT", slog.String("error", err.Error()))
		return
	}

	return
}

func (a *AuthServiceImpl) LogIn(ctx context.Context, key, password string, serviceId int64) (token string, err error) {
	return "", nil
}

func (a *AuthServiceImpl) LogOut(ctx context.Context, token string) error {
	return fmt.Errorf("not implemented")
}

func initLogin(login, email, phone string) (string, error) {
	switch {
	case login != "":
		return login, nil
	case email != "":
		return email, nil
	case phone != "":
		return phone, nil
	default:
		return "", ErrInvalidCredentials
	}
}

func nilIfEmpty(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

func (a *AuthServiceImpl) userByKey(ctx context.Context, logger *slog.Logger, key string) (user *models.SavedUser, err error) {
	validate := validator.New()
	logger.Info("validating key")

	switch {
	case validate.Var(key, "email") == nil:
		user, err = a.db.UserByEmail(ctx, key)
		logger = logger.With(slog.String("key_type", "email"))
	case validate.Var(key, "e164") == nil:
		user, err = a.db.UserByPhone(ctx, key)
		logger = logger.With(slog.String("key_type", "phone"))
	default:
		user, err = a.db.UserByLogin(ctx, key)
		logger = logger.With(slog.String("key_type", "login"))
	}
	if err != nil {
		logger.Info("cannot get user", slog.String("error", err.Error()))
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return
	}
	return
}
