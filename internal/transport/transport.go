package transport

import (
	"context"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/models"
)

type AuthService interface {
	Register(ctx context.Context, creds *models.Credentials) (*models.LoginResponse, ero.Error)
	Login(ctx context.Context, creds *models.Credentials) (*models.LoginResponse, ero.Error)
	Logout(ctx context.Context, refreshToken string) ero.Error
	Relogin(ctx context.Context, refreshToken string) (*models.LoginResponse, ero.Error)
	Check(ctx context.Context, superAccessToken string) ero.Error
	SetProfile(ctx context.Context, anyAccessToken string, user *models.User) ero.Error
	GetProfile(ctx context.Context, superAccessToken string) (*models.Profile, ero.Error)
	GetApps(ctx context.Context, superAccessToken string) ([]models.App, ero.Error)
	SetPassword(ctx context.Context, superAccessToken string, password string) ero.Error
	Delete(ctx context.Context, superAccessToken string, password string) ero.Error
	DeleteApp(ctx context.Context, superAccessToken string, appId int64) ero.Error
	Recover(ctx context.Context, creds *models.Credentials) ero.Error
}

type AppService interface {
	Login(ctx context.Context, data *models.AppCredentials) (*models.LoginResponse, ero.Error)
	Logout(ctx context.Context, refreshToken string) ero.Error
	Relogin(ctx context.Context, refreshToken string) (*models.LoginResponse, ero.Error)
	Check(ctx context.Context, accessToken string) ero.Error
	SetProfile(ctx context.Context, anyAccessToken string, user *models.User) ero.Error
	GetProfile(ctx context.Context, accessToken string) (*models.Profile, ero.Error)
	GetSessions(ctx context.Context, accessToken string) ([]models.SessionInfo, ero.Error)
	Delete(ctx context.Context, accessToken string) ero.Error
	Recover(ctx context.Context, data *models.AppCredentials) ero.Error
}
