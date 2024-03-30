package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooLong    = errors.New("password is longer than 72 bytes")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (a *AuthServiceImpl) SignUp(ctx context.Context, optLogin, optEmail, optPhone, appToken, password string) (token string, err error) {
	const op = "auth.AuthServiceImpl.SignUp"
	logger := slog.With(slog.String("op", op))

	optLogin, err = initLogin(optLogin, optEmail, optPhone)
	if err != nil {
		return
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrPasswordTooLong
	}

	user := &models.User{
		Login:    optLogin,
		Email:    optEmail,
		Phone:    optPhone,
		Password: string(passHash),
	}

	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		logger.Error("Could not sign up because of user validation errors", slog.Any("errors", errs.Error()))
		return "", ErrInvalidCredentials
	}

	saved, err := a.db.SaveUser(ctx, user)
	if err != nil {
		return
	}
	// TODO: SaveSignup here

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    saved.Id,
		"login": saved.Login,
		"email": saved.Email,
		"phone": saved.Phone,
		"exp":   time.Now().Add(a.tokenTTL).Unix(),
	})

	token, err = jwtToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return
	}
	return
}

func (a *AuthServiceImpl) LogIn(ctx context.Context, optLogin, optEmail, optPhone, appToken, password string) (token string, err error) {
	return "", fmt.Errorf("not implemented")
}

func (a *AuthServiceImpl) LogOut(ctx context.Context, token string) error {
	return fmt.Errorf("not implemented")
}

func initLogin(login, email, phone string) (string, error) {
	if login != "" {
		return login, nil
	}

	switch {
	case email != "":
		return email, nil
	case phone != "":
		return phone, nil
	default:
		return "", fmt.Errorf("at least one of login, email or phone must be set")
	}
}
