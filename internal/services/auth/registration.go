package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	goos "os"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/Onnywrite/grpc-auth/internal/storage"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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
	ErrInvalidIP          = errors.New("invalid IP")
	ErrUnauthorized       = errors.New("unauthorized")
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

func (a *AuthServiceImpl) LogIn(ctx context.Context, key, password string, serviceId int64, browser, os, ip string) (*models.Tokens, error) {
	const op = "auth.AuthServiceImpl.LogIn"
	log := a.log.With(slog.String("op", op), slog.String("key", key),
		slog.Int64("service_id", serviceId), slog.String("ip", ip), slog.String("os", os), slog.String("browser", browser))

	user, err := a.recognize(ctx, key)
	if err != nil {
		log.Error("recognition failed")
		return nil, err
	}
	log = log.With(slog.Int64("user_id", user.Id))

	su, err := a.db.Signup(ctx, user.Id, serviceId)
	if err != nil {
		log.Error("signup not found")
	}

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
	// TODO:
	if err != nil {
		log.Error("could not save session")
		return nil, ErrInternal
	}

	refreshTkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"session_uuid": session.UUID.String(),
		"exp":          time.Now().Add(a.refreshTokenTTL).Unix(),
	})

	refresh, err := refreshTkn.SignedString(goos.Getenv(tokenEnv))
	if err != nil {
		log.Error("failed to sign access token", slog.String("error", err.Error()))
		return nil, ErrInternal
	}

	accessTkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         su.UserId,
		"login":      user.Login,
		"service_id": su.ServiceId,
		// probably custom type like models.Role
		"roles": []string{},
		"exp":   time.Now().Add(a.tokenTTL).Unix(),
	})

	access, err := accessTkn.SignedString(goos.Getenv(tokenEnv))
	if err != nil {
		log.Error("failed to sign access token", slog.String("error", err.Error()))
		return nil, ErrInternal
	}

	return &models.Tokens{
		Refresh: refresh,
		Access:  access,
	}, nil
}

func (a *AuthServiceImpl) LogOut(ctx context.Context, refresh string) error {
	const op = "auth.AuthServiceImpl.LogOut"
	log := a.log.With(slog.String("op", op))

	token, err := a.processRefreshToken(refresh)
	if err != nil {
		log.Error("could not process token", slog.String("token", refresh), slog.String("error", err.Error()))
		return ErrUnauthorized
	}
	log = log.With(slog.String("session_uuid", token.SessionUUID.String()))

	err = a.db.DeleteSession(ctx, token.SessionUUID)
	if err != nil {
		log.Error("could not delete session", slog.String("error", err.Error()))
		// TODO:
		return errors.New("TODO")
	}

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

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrTokenExpired            = errors.New("token has expired")
)

func (a *AuthServiceImpl) processAccessToken(tkn string) (*models.AccessToken, error) {
	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return []byte(goos.Getenv(tokenEnv)), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := claims["exp"].(float64)

		if float64(time.Now().Unix()) > exp {
			return nil, ErrTokenExpired
		}

		id, ok := claims["id"].(float64)
		if !ok {
			return nil, fmt.Errorf("could not convert 'id' to int64")
		}
		login, ok := claims["login"].(string)
		if !ok {
			return nil, fmt.Errorf("could not convert 'login' to string")
		}
		serviceId, ok := claims["service_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("could not convert 'service_id' to int64")
		}
		roles, ok := claims["roles"].([]string)
		if !ok {
			return nil, fmt.Errorf("could not convert 'roles' to []string")
		}

		token := &models.AccessToken{
			Id:        int64(id),
			Login:     login,
			ServiceId: int64(serviceId),
			Roles:     roles,
			Exp:       int64(exp),
		}

		return token, nil
	}

	return nil, err
}

func (a *AuthServiceImpl) processRefreshToken(tkn string) (*models.RefreshToken, error) {
	const op = "auth.AuthServiceImpl.processRefreshToken"
	log := a.log.With(slog.String("op", op))

	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		return []byte(goos.Getenv(tokenEnv)), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := claims["exp"].(float64)

		if float64(time.Now().Unix()) > exp {
			return nil, ErrTokenExpired
		}

		sessionUUIDStr, ok := claims["session_uuid"].(string)
		if !ok {
			log.Error("could not convert 'session_uuid' to uuid.UUID")
			return nil, ErrInternal
		}
		sessionUUID, err := uuid.Parse(sessionUUIDStr)
		if err != nil {
			log.Error("could not convert 'session_uuid' to uuid.UUID",
				slog.String("session_uuid_str", sessionUUIDStr))
			return nil, ErrInternal
		}

		token := &models.RefreshToken{
			SessionUUID: sessionUUID,
			Exp:         int64(exp),
		}

		return token, nil
	}

	return nil, err
}
