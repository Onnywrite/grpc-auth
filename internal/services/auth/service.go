package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/models"
)

type Storage interface {
	SaveUser(ctx context.Context, user *models.User) (*models.SavedUser, error)
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
