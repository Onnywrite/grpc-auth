package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/models"
)

type Storage interface {
	SaveUser(ctx context.Context, user *models.User) (*models.SavedUser, error)
	SaveSignup(ctx context.Context, signup models.Signup) error

	UserById(ctx context.Context, id int64) (*models.SavedUser, error)
	UserByLogin(ctx context.Context, login string) (*models.SavedUser, error)
	UserByEmail(ctx context.Context, email string) (*models.SavedUser, error)
	UserByPhone(ctx context.Context, phone string) (*models.SavedUser, error)

	Signup(ctx context.Context, userId, serviceId int64) (*models.SavedSignup, error)
}

type AuthServiceImpl struct {
	log      *slog.Logger
	db       Storage
	tokenTTL time.Duration
}

func New(logger *slog.Logger, db Storage, tokenTTL time.Duration) *AuthServiceImpl {
	return &AuthServiceImpl{
		log:      logger,
		db:       db,
		tokenTTL: tokenTTL,
	}
}
