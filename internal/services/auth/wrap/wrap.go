package wrap

import (
	"context"
	"log/slog"

	"github.com/Onnywrite/grpc-auth/internal/models"
)

type Storage interface {
	SaveUser(ctx context.Context, user *models.User) (*models.SavedUser, error)
	UserById(ctx context.Context, id int64) (*models.SavedUser, error)
	UserByLogin(ctx context.Context, login string) (*models.SavedUser, error)
	UserByEmail(ctx context.Context, email string) (*models.SavedUser, error)
	UserByPhone(ctx context.Context, phone string) (*models.SavedUser, error)

	SaveSignup(ctx context.Context, signup models.Signup) (*models.SavedSignup, error)
	SignupById(ctx context.Context, id int64) (*models.SavedSignup, error)
	SignupByServiceAndUser(ctx context.Context, serviceId, userId int64) (*models.SavedSignup, error)

	SaveSession(ctx context.Context, session *models.Session) (*models.SavedSession, error)
	SessionByUuid(ctx context.Context, uuid string) (*models.SavedSession, error)
	SessionByInfo(ctx context.Context, signupId int64, info models.SessionInfo) (*models.SavedSession, error)
	TerminateSession(ctx context.Context, uuid string) error
	ReviveSession(ctx context.Context, uuid string) error
	DeleteSession(ctx context.Context, uuid string) error
}

type Wrapper struct {
	log *slog.Logger
	db  Storage
}

func New(logger *slog.Logger, db Storage) *Wrapper {
	return &Wrapper{
		log: logger,
		db:  db,
	}
}
