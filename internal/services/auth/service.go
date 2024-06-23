package auth

import (
	"context"
	"log/slog"

	"github.com/Onnywrite/grpc-auth/internal/config"
	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/models"
)

type Storage interface {
	SaveUser(ctx context.Context, user *models.User) (*models.SavedUser, ero.Error)
	UserById(ctx context.Context, id int64) (*models.SavedUser, ero.Error)
	UserByNickname(ctx context.Context, nickname string) (*models.SavedUser, ero.Error)
	UserByEmail(ctx context.Context, email string) (*models.SavedUser, ero.Error)
	UserByPhone(ctx context.Context, phone string) (*models.SavedUser, ero.Error)

	SaveSignup(ctx context.Context, signup models.Signup) (*models.SavedSignup, ero.Error)
	SignupByServiceAndUser(ctx context.Context, serviceId, userId int64) (*models.SavedSignup, ero.Error)

	SaveSession(ctx context.Context, session *models.Session) (*models.SavedSession, ero.Error)
	SessionByUuid(ctx context.Context, uuid string) (*models.SavedSession, ero.Error)
	SessionByInfo(ctx context.Context, serviceId, userId int64, info models.SessionInfo) (*models.SavedSession, ero.Error)
	DeleteSession(ctx context.Context, uuid string) ero.Error
}

type AuthService struct {
	log           *slog.Logger
	db            Storage
	accessConfig  *config.TokenConfig
	refreshConfig *config.TokenConfig
}

func New(logger *slog.Logger, db Storage, accessConfig, refreshConfig *config.TokenConfig) *AuthService {
	return &AuthService{
		log:           logger,
		db:            NewWrapper(logger, db),
		accessConfig:  accessConfig,
		refreshConfig: refreshConfig,
	}
}
