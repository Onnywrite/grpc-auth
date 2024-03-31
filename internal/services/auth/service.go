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

	UserById(ctx context.Context, id int64) (u *models.SavedUser, err error)
	UserByLogin(ctx context.Context, login string) (u *models.SavedUser, err error)
	UserByEmail(ctx context.Context, email string) (u *models.SavedUser, err error)
	UserByPhone(ctx context.Context, phone string) (u *models.SavedUser, err error)
	// change models.SavedUser to models.SavedSignup
	//SignupByLogin(ctx context.Context, login string, serviceId int64) (u *models.SavedUser, err error)
	//SignupByEmail(ctx context.Context, email string, serviceId int64) (u *models.SavedUser, err error)
	//SignupByPhone(ctx context.Context, phone string, serviceId int64) (u *models.SavedUser, err error)
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
