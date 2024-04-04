package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/google/uuid"
)

type Storage interface {
	SaveUser(ctx context.Context, user *models.User) (*models.SavedUser, error)
	SaveSignup(ctx context.Context, signup models.Signup) error
	SaveSession(ctx context.Context, session *models.Session) (*models.SavedSession, error)

	UserById(ctx context.Context, id int64) (*models.SavedUser, error)
	UserByLogin(ctx context.Context, login string) (*models.SavedUser, error)
	UserByEmail(ctx context.Context, email string) (*models.SavedUser, error)
	UserByPhone(ctx context.Context, phone string) (*models.SavedUser, error)

	Signup(ctx context.Context, userId, serviceId int64) (*models.SavedSignup, error)

	TerminateSession(ctx context.Context, uuid uuid.UUID) error
}

type AuthServiceImpl struct {
	log                       *slog.Logger
	db                        Storage
	tokenTTL, refreshTokenTTL time.Duration
}

func New(logger *slog.Logger, db Storage, tokenTTL, refreshTokenTTL time.Duration) *AuthServiceImpl {
	return &AuthServiceImpl{
		log:      logger,
		db:       db,
		tokenTTL: tokenTTL,
	}
}
