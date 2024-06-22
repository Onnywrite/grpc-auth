package auth

import (
	"context"
	"log/slog"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Register(ctx context.Context, creds *models.Credentials) (*models.LoginResponse, ero.Error) {
	const op = "auth.AuthService.Register"
	log := s.log.With(slog.String("op", op), slog.String("nickname", creds.Nickname))

	if erro := s.validateCredentials(ctx, log, creds); erro != nil {
		return nil, erro
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", slog.String("error", err.Error()))
		return nil, ero.NewServer(op)
	}
	creds.Password = string(hashed)

	saved, erro := s.db.SaveUser(ctx, &creds.User)
	if erro != nil {
		log.Error("failed to save user", slog.String("error", erro.Error()))
		return nil, ero.NewServer(op)
	}
	log.Info("saved user")

	return s.openSession(ctx, log, saved, creds.Info)
}

func (s *AuthService) Login(ctx context.Context, creds *models.Credentials) (*models.LoginResponse, ero.Error) {
	const op = "AuthService.Nickname"
	log := s.log.With(slog.String("op", op), slog.String("nickname", creds.Nickname))

	if erro := s.validateCredentials(ctx, log, creds); erro != nil {
		return nil, erro
	}

	user, erro := s.db.UserByNickname(ctx, creds.Nickname)
	if erro != nil {
		log.Error("failed to get user by nickname", slog.String("error", erro.Error()))
		return nil, erro
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		log.Error("hash and passwword mismatch", slog.String("error", err.Error()))
		return nil, ero.NewClient("invalid credentials")
	}

	return s.openSession(ctx, log, user, creds.Info)
}

func (*AuthService) Logout(ctx context.Context, refreshToken string) ero.Error {
	return nil
}
func (*AuthService) Relogin(ctx context.Context, refreshToken string) (*models.LoginResponse, ero.Error) {
	return nil, nil
}
func (*AuthService) Check(ctx context.Context, superAccessToken string) ero.Error {
	return nil
}
func (*AuthService) SetProfile(ctx context.Context, anyAccessToken string, user *models.User) ero.Error {
	return nil
}
func (*AuthService) GetProfile(ctx context.Context, superAccessToken string) (*models.Profile, ero.Error) {
	return nil, nil
}
func (*AuthService) GetApps(ctx context.Context, superAccessToken string) ([]models.App, ero.Error) {
	return nil, nil
}
func (*AuthService) SetPassword(ctx context.Context, superAccessToken string, password string) ero.Error {
	return nil
}
func (*AuthService) Delete(ctx context.Context, superAccessToken string, password string) ero.Error {
	return nil
}
func (*AuthService) DeleteApp(ctx context.Context, superAccessToken string, appId int64) ero.Error {
	return nil
}
func (*AuthService) Recover(ctx context.Context, creds *models.Credentials) ero.Error {
	return nil
}
