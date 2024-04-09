package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/Onnywrite/grpc-auth/internal/models"
)

type Storage interface {
	SaveUser(ctx context.Context, user *models.User) (*models.SavedUser, error)
	User(ctx context.Context, user models.UserIdentifier) (*models.SavedUser, error)

	SaveSignup(ctx context.Context, signup models.Signup) (*models.SavedSignup, error)
	Signup(ctx context.Context, userId, serviceId int64) (*models.SavedSignup, error)

	SaveSession(ctx context.Context, session *models.Session) (*models.SavedSession, error)
	SessionById(ctx context.Context, uuid string) (*models.SavedSession, error)
	Session(ctx context.Context, session *models.Session) (*models.SavedSession, error)
	TerminateSession(ctx context.Context, uuid string) error
	// ReviveSession(ctx context.Context, uuid string) error
	DeleteSession(ctx context.Context, uuid string) error
}

type AuthService struct {
	log                       *slog.Logger
	db                        Storage
	tokenTTL, refreshTokenTTL time.Duration
}

func New(logger *slog.Logger, db Storage, tokenTTL, refreshTokenTTL time.Duration) *AuthService {
	return &AuthService{
		log:      logger,
		db:       db,
		tokenTTL: tokenTTL,
	}
}
