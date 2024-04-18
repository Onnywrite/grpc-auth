package wrap

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Onnywrite/grpc-auth/internal/models"
	auth "github.com/Onnywrite/grpc-auth/internal/services/auth/common"
	storage "github.com/Onnywrite/grpc-auth/internal/storage/common"
)

type Storage interface {
	SaveUser(ctx context.Context, user *models.User) (*models.SavedUser, error)
	UserById(ctx context.Context, id int64) (*models.SavedUser, error)
	UserByLogin(ctx context.Context, login string) (*models.SavedUser, error)
	UserByEmail(ctx context.Context, email string) (*models.SavedUser, error)
	UserByPhone(ctx context.Context, phone string) (*models.SavedUser, error)

	SaveSignup(ctx context.Context, signup models.Signup) (*models.SavedSignup, error)
	SignupByServiceAndUser(ctx context.Context, serviceId, userId int64) (*models.SavedSignup, error)

	SaveSession(ctx context.Context, session *models.Session) (*models.SavedSession, error)
	SessionByUuid(ctx context.Context, uuid string) (*models.SavedSession, error)
	SessionByInfo(ctx context.Context, serviceId, userId int64, info models.SessionInfo) (*models.SavedSession, error)
	TerminateSession(ctx context.Context, uuid string) error
	ReviveSession(ctx context.Context, uuid string) error
	DeleteSession(ctx context.Context, uuid string) error
}

type Wrapper struct {
	log *slog.Logger
	Storage
}

func New(logger *slog.Logger, db Storage) *Wrapper {
	return &Wrapper{
		log:     logger,
		Storage: db,
	}
}

// Throws:
//
//	ErrSignupNotExists
//	ErrInternal
func (w *Wrapper) SaveSession(ctx context.Context, session *models.Session) (*models.SavedSession, func() error, error) {
	const op = "wrap.Wrapper.SaveSession"
	log := w.log.With(slog.String("op", op), slog.Int64("service_id", session.ServiceId), slog.Int64("user_id", session.UserId))

	saved, err := w.Storage.SaveSession(ctx, session)
	if errors.Is(err, storage.ErrFKConstraint) {
		log.Error("user has never signed up")
		return nil, nil, auth.ErrSignupNotExists
	}
	if errors.Is(err, storage.ErrUniqueConstraint) {
		saved, err = w.SessionByInfo(ctx, session.ServiceId, session.UserId, session.Info)
		if errors.Is(err, auth.ErrSessionTerminated) {
			err = w.Storage.ReviveSession(ctx, saved.UUID)
			if err != nil {
				log.Error("error reviving session", slog.String("error", err.Error()))
				return nil, nil, auth.ErrInternal
			}
			log.Info("session revived", slog.String("uuid", saved.UUID))
			return saved, func() error {
				return w.TerminateSession(context.Background(), saved.UUID)
			}, nil
		}
		if err != nil {
			log.Error("error getting session", slog.String("error", err.Error()))
			return nil, nil, auth.ErrInternal
		}
	}

	if err != nil {
		log.Error("error saving session", slog.String("error", err.Error()))
		return nil, nil, auth.ErrInternal
	}

	log.Info("session saved", slog.String("uuid", saved.UUID))

	return saved, func() error {
		return w.DeleteSession(ctx, saved.UUID)
	}, nil
}

// Throws:
//
//	ErrSessionNotExists
//	ErrSessionTerminated
//	ErrInternal
func (w *Wrapper) SessionByUuid(ctx context.Context, uuid string) (*models.SavedSession, error) {
	return w.session(ctx, func(ctx context.Context, keys ...any) (*models.SavedSession, error) {
		return w.Storage.SessionByUuid(ctx, keys[0].(string))
	}, uuid)
}

// Throws:
//
//	ErrSessionNotExists
//	ErrSessionTerminated
//	ErrInternal
func (w *Wrapper) SessionByInfo(ctx context.Context, serviceId, userId int64, info models.SessionInfo) (*models.SavedSession, error) {
	return w.session(ctx, func(ctx context.Context, keys ...any) (*models.SavedSession, error) {
		return w.Storage.SessionByInfo(ctx, keys[0].(int64), keys[1].(int64), keys[2].(models.SessionInfo))
	}, serviceId, userId, info)
}

type getSessionFn func(ctx context.Context, keys ...any) (*models.SavedSession, error)

// Throws:
//
//	ErrSessionNotExists
//	ErrSessionTerminated
//	ErrInternal
func (w *Wrapper) session(ctx context.Context, get getSessionFn, keys ...any) (*models.SavedSession, error) {
	const op = "w.Wrapper.session"
	log := w.log.With(slog.String("op", op))

	session, err := get(ctx, keys...)
	if errors.Is(err, storage.ErrEmptyResult) {
		log.Error("session not found")
		return nil, auth.ErrSessionNotExists
	}
	if err != nil {
		log.Error("error getting session", slog.String("error", err.Error()))
		return nil, auth.ErrInternal
	}

	if session.IsTerminated() {
		log.Error("session terminated")
		return nil, auth.ErrSessionTerminated
	}

	log.Info("got session", slog.String("uuid", session.UUID))

	return session, nil
}

// Throws:
//
//	ErrSignedOut
//	ErrUserAlreadyRegistered
//	ErrServiceNotExists
//	ErrInternal
func (w *Wrapper) SaveSignup(ctx context.Context, signup models.Signup) (*models.SavedSignup, error) {
	const op = "wrap.Wrapper.SaveSignup"
	log := w.log.With(slog.String("op", op), slog.Int64("service_id", signup.ServiceId), slog.Int64("user_id", signup.UserId))

	su, err := w.Storage.SaveSignup(ctx, signup)
	if errors.Is(err, storage.ErrUniqueConstraint) {
		su, err = w.Storage.SignupByServiceAndUser(ctx, signup.UserId, signup.ServiceId)
		if err != nil {
			log.Error("error getting signup", slog.String("error", err.Error()))
			return nil, auth.ErrInternal
		}
		if su.IsDeleted() {
			log.Error("signed out")
			return nil, auth.ErrSignedOut
		}
		if su.IsBanned() {
			log.Error("signup banned")
			return nil, auth.ErrSignupBanned
		}

		return nil, auth.ErrAlreadySignedUp
	}

	if errors.Is(err, storage.ErrFKConstraint) {
		log.Error("service does not exist", slog.String("error", err.Error()))
		return nil, auth.ErrServiceNotExists
	}

	if err != nil {
		log.Error("saving error", slog.String("error", err.Error()))
		return nil, auth.ErrInternal
	}
	log.Info("saved signup")

	return su, nil
}

// Throws:
//
//	ErrSignupNotExists
//	ErrSignedOut
//	ErrSignupBanned
//	ErrInternal
func (w *Wrapper) SignupByServiceAndUser(ctx context.Context, serviceId, userId int64) (*models.SavedSignup, error) {
	return w.signup(ctx, func(ctx context.Context, keys ...any) (*models.SavedSignup, error) {
		return w.Storage.SignupByServiceAndUser(ctx, keys[0].(int64), keys[1].(int64))
	}, serviceId, userId)
}

type getSignupFn func(ctx context.Context, keys ...any) (*models.SavedSignup, error)

// Throws:
//
//	ErrSignupNotExists
//	ErrSignedOut
//	ErrSignupBanned
//	ErrInternal
func (w *Wrapper) signup(ctx context.Context, get getSignupFn, keys ...any) (*models.SavedSignup, error) {
	const op = "wrap.Wrapper.signupBy"
	log := w.log.With(slog.String("op", op), slog.Any("keys", keys))

	su, err := get(ctx, keys...)
	if errors.Is(err, storage.ErrEmptyResult) {
		log.Error("signup not found")
		return nil, auth.ErrSignupNotExists
	}
	if err != nil {
		log.Error("getting error", slog.String("error", err.Error()))
		return nil, auth.ErrInternal
	}

	if su.IsDeleted() {
		log.Error("signed out")
		return nil, auth.ErrSignedOut
	}

	if su.IsBanned() {
		log.Error("signup banned")
		return nil, auth.ErrSignupBanned
	}

	log.Info("got signup")

	return su, nil
}

// Throws;
//
//	ErrUserAlreadyRegistered
//	ErrUserDeleted
//	ErrInternal in any unexpected situation
func (w *Wrapper) SaveUser(ctx context.Context, user *models.User) (*models.SavedUser, error) {
	const op = "wrap.Wrapper.SaveUser"
	log := w.log.With(slog.String("op", op))

	u, err := w.Storage.SaveUser(ctx, user)
	if errors.Is(err, storage.ErrUniqueConstraint) {
		u, err = w.UserByLogin(ctx, *user.Login)
		if errors.Is(err, auth.ErrUserDeleted) {
			log.Error("user deleted", slog.Int64("id", u.Id), slog.String("error", err.Error()))
			return nil, err
		}
		log.Error("user already exists", slog.Int64("id", u.Id), slog.String("error", err.Error()))
		return nil, auth.ErrUserAlreadyRegistered
	}
	if err != nil {
		log.Error("saving error", slog.String("error", err.Error()))
		return nil, auth.ErrInternal
	}
	log.Info("saved user", slog.Int64("id", u.Id))

	return u, nil
}

// Throws;
//
//	ErrInvalidCredentials if nothing found
//	ErrUserDeleted
//	ErrInternal in any unexpected situation
func (w *Wrapper) UserByLogin(ctx context.Context, login string) (*models.SavedUser, error) {
	return w.user(ctx, func(ctx context.Context, key any) (*models.SavedUser, error) {
		return w.Storage.UserByLogin(ctx, key.(string))
	}, login, "login")
}

// Throws;
//
//	ErrInvalidCredentials if nothing found
//	ErrUserDeleted
//	ErrInternal in any unexpected situation
func (w *Wrapper) UserByEmail(ctx context.Context, email string) (*models.SavedUser, error) {
	return w.user(ctx, func(ctx context.Context, key any) (*models.SavedUser, error) {
		return w.Storage.UserByEmail(ctx, key.(string))
	}, email, "email")
}

// Throws;
//
//	ErrInvalidCredentials if nothing found
//	ErrUserDeleted
//	ErrInternal in any unexpected situation
func (w *Wrapper) UserByPhone(ctx context.Context, phone string) (*models.SavedUser, error) {
	return w.user(ctx, func(ctx context.Context, key any) (*models.SavedUser, error) {
		return w.Storage.UserByLogin(ctx, key.(string))
	}, phone, "phone")
}

// Throws;
//
//	ErrInvalidCredentials if nothing found
//	ErrUserDeleted
//	ErrInternal in any unexpected situation
func (w *Wrapper) UserById(ctx context.Context, id int64) (*models.SavedUser, error) {
	return w.user(ctx, func(ctx context.Context, key any) (*models.SavedUser, error) {
		return w.Storage.UserById(ctx, key.(int64))
	}, id, "id")
}

type getUserFn func(ctx context.Context, key any) (*models.SavedUser, error)

// Throws;
//
//	ErrInvalidCredentials if nothing found
//	ErrUserDeleted
//	ErrInternal in any unexpected situation
func (w *Wrapper) user(ctx context.Context, get getUserFn, key any, keyType string) (*models.SavedUser, error) {
	const op = "wrap.Wrapper.user"
	log := w.log.With(slog.String("op", op), slog.Any(keyType, key))

	saved, err := get(ctx, key)
	if errors.Is(err, storage.ErrEmptyResult) {
		log.Error("no user found")
		return nil, auth.ErrInvalidCredentials
	}
	if err != nil {
		log.Error("getting error", slog.String("error", err.Error()))
		return nil, auth.ErrInternal
	}

	if saved.IsDeleted() {
		log.Error("user deleted")
		return saved, auth.ErrUserDeleted
	}

	return saved, nil
}
