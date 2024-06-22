package auth

import (
	"context"
	"log/slog"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/models"
	storage "github.com/Onnywrite/grpc-auth/internal/storage/common"
)

type Wrapper struct {
	log *slog.Logger
	Storage
}

func NewWrapper(logger *slog.Logger, db Storage) *Wrapper {
	return &Wrapper{
		log:     logger,
		Storage: db,
	}
}

// Throws:
//
//	ErrSignupNotExists
//	ErrSessionAlreadyOpened
//	Internal
func (w *Wrapper) SaveSession(ctx context.Context, session *models.Session) (*models.SavedSession, ero.Error) {
	const op = "wrap.Wrapper.SaveSession"
	log := w.log.With(slog.String("op", op), slog.Int64("service_id", session.ServiceId), slog.Int64("user_id", session.UserId))

	saved, err := w.Storage.SaveSession(ctx, session)
	switch {
	case ero.Has(err, storage.ErrFKConstraint):
		log.Error("user has never signed up")
		return nil, ero.NewClient(ErrSignupNotExists)
	case ero.Has(err, storage.ErrUniqueConstraint):
		log.Error("session already exists")
		return nil, ero.NewClient(ErrSessionAlreadyOpened)
	case err != nil:
		log.Error("error saving session", slog.String("error", err.Error()))
		return nil, ero.InternalFrom(op, err)
	}

	log.Info("session saved", slog.String("uuid", saved.UUID))

	return saved, nil
}

// Throws:
//
//	ErrSessionNotExists
//	Internal
func (w *Wrapper) SessionByUuid(ctx context.Context, uuid string) (*models.SavedSession, ero.Error) {
	return w.session(ctx, func(ctx context.Context, keys ...any) (*models.SavedSession, ero.Error) {
		return w.Storage.SessionByUuid(ctx, keys[0].(string))
	}, uuid)
}

// Throws:
//
//	ErrSessionNotExists
//	Internal
func (w *Wrapper) SessionByInfo(ctx context.Context, serviceId, userId int64, info models.SessionInfo) (*models.SavedSession, ero.Error) {
	return w.session(ctx, func(ctx context.Context, keys ...any) (*models.SavedSession, ero.Error) {
		return w.Storage.SessionByInfo(ctx, keys[0].(int64), keys[1].(int64), keys[2].(models.SessionInfo))
	}, serviceId, userId, info)
}

type getSessionFn func(ctx context.Context, keys ...any) (*models.SavedSession, ero.Error)

// Throws:
//
//	ErrSessionNotExists
//	Internal
func (w *Wrapper) session(ctx context.Context, get getSessionFn, keys ...any) (*models.SavedSession, ero.Error) {
	const op = "w.Wrapper.session"
	log := w.log.With(slog.String("op", op))

	session, err := get(ctx, keys...)
	switch {
	case ero.Has(err, storage.ErrEmptyResult):
		log.Error("session not found")
		return nil, ero.NewClient(ErrSessionNotExists)
	case err != nil:
		log.Error("error getting session", slog.String("error", err.Error()))
		return nil, ero.InternalFrom(op, err)
	}

	log.Info("got session", slog.String("uuid", session.UUID))

	return session, nil
}

// Throws:
//
//	ErrServiceNotExists
//	ErrSignedOut
//	ErrSignupBanned
//	ErrAlreadySignedUp
//	Internal
func (w *Wrapper) SaveSignup(ctx context.Context, signup models.Signup) (*models.SavedSignup, ero.Error) {
	const op = "wrap.Wrapper.SaveSignup"
	log := w.log.With(slog.String("op", op), slog.Int64("service_id", signup.ServiceId), slog.Int64("user_id", signup.UserId))

	su, err := w.Storage.SaveSignup(ctx, signup)
	switch {
	case ero.Has(err, storage.ErrFKConstraint):
		log.Error("service does not exist")
		return nil, ero.NewClient(ErrServiceNotExists)
	case ero.Has(err, storage.ErrUniqueConstraint):
		log.Error("user already signed up, checking details")
		su, err = w.Storage.SignupByServiceAndUser(ctx, signup.ServiceId, signup.UserId)
		if err != nil {
			log.Error("error getting signup", slog.String("error", err.Error()))
			return nil, ero.InternalFrom(op, err)
		}
		if su.IsDeleted() {
			log.Error("signed out")
			return nil, ero.NewClient(ErrSignedOut)
		}
		if su.IsBanned() {
			log.Error("signup banned")
			return nil, ero.NewClient(ErrSignupBanned)
		}
		return nil, ero.NewClient(ErrAlreadySignedUp)
	case err != nil:
		log.Error("saving error", slog.String("error", err.Error()))
		return nil, ero.InternalFrom(op, err)
	}

	log.Info("saved signup")

	return su, nil
}

// Throws:
//
//	ErrSignupNotExists
//	ErrSignedOut
//	ErrSignupBanned
//	Internal
func (w *Wrapper) SignupByServiceAndUser(ctx context.Context, serviceId, userId int64) (*models.SavedSignup, ero.Error) {
	return w.signup(ctx, func(ctx context.Context, keys ...any) (*models.SavedSignup, ero.Error) {
		return w.Storage.SignupByServiceAndUser(ctx, keys[0].(int64), keys[1].(int64))
	}, serviceId, userId)
}

type getSignupFn func(ctx context.Context, keys ...any) (*models.SavedSignup, ero.Error)

// Throws:
//
//	ErrSignupNotExists
//	ErrSignedOut
//	ErrSignupBanned
//	Internal
func (w *Wrapper) signup(ctx context.Context, get getSignupFn, keys ...any) (*models.SavedSignup, ero.Error) {
	const op = "wrap.Wrapper.signupBy"
	log := w.log.With(slog.String("op", op), slog.Any("keys", keys))

	su, err := get(ctx, keys...)
	switch {
	case ero.Has(err, storage.ErrEmptyResult):
		log.Error("signup not found")
		return nil, ero.NewClient(ErrSignupNotExists)
	case err != nil:
		log.Error("cannot get signup", slog.String("error", err.Error()))
		return nil, ero.InternalFrom(op, err)
	}

	if su.IsDeleted() {
		log.Error("signed out")
		return su, ero.NewClient(ErrSignedOut)
	}

	if su.IsBanned() {
		log.Error("signup banned")
		return su, ero.NewClient(ErrSignupBanned)
	}

	log.Info("got signup")

	return su, nil
}

// Throws;
//
//	ErrUserAlreadyRegistered
//	ErrUserDeleted
//	Internal
func (w *Wrapper) SaveUser(ctx context.Context, user *models.User) (*models.SavedUser, ero.Error) {
	const op = "wrap.Wrapper.SaveUser"
	log := w.log.With(slog.String("op", op))

	u, err := w.Storage.SaveUser(ctx, user)
	switch {
	case ero.Has(err, storage.ErrUniqueConstraint):
		u, err = w.UserByNickname(ctx, user.Nickname)
		if err.Has(ErrUserDeleted) {
			log.Error("user deleted", slog.Int64("id", u.Id), slog.String("error", err.Error()))
			return nil, err
		}
		log.Error("user already exists", slog.Int64("id", u.Id))
		return nil, ero.NewClient(ErrUserAlreadyRegistered)
	case err != nil:
		log.Error("cannot save user", slog.String("error", err.Error()))
		return nil, ero.InternalFrom(op, err)
	}

	log.Info("saved user", slog.Int64("id", u.Id))

	return u, nil
}

// Throws:
//
//	ErrInvalidCredentials
//	ErrUserDeleted
//	Internal
func (w *Wrapper) UserByNickname(ctx context.Context, nickname string) (*models.SavedUser, ero.Error) {
	return w.user(ctx, func(ctx context.Context, key any) (*models.SavedUser, ero.Error) {
		return w.Storage.UserByNickname(ctx, key.(string))
	}, nickname, "nickname")
}

// Throws:
//
//	ErrInvalidCredentials
//	ErrUserDeleted
//	Internal
func (w *Wrapper) UserByEmail(ctx context.Context, email string) (*models.SavedUser, ero.Error) {
	return w.user(ctx, func(ctx context.Context, key any) (*models.SavedUser, ero.Error) {
		return w.Storage.UserByEmail(ctx, key.(string))
	}, email, "email")
}

// Throws:
//
//	ErrInvalidCredentials
//	ErrUserDeleted
//	Internal
func (w *Wrapper) UserByPhone(ctx context.Context, phone string) (*models.SavedUser, ero.Error) {
	return w.user(ctx, func(ctx context.Context, key any) (*models.SavedUser, ero.Error) {
		return w.Storage.UserByPhone(ctx, key.(string))
	}, phone, "phone")
}

// Throws;
//
//	ErrInvalidCredentials
//	ErrUserDeleted
//	Internal
func (w *Wrapper) UserById(ctx context.Context, id int64) (*models.SavedUser, ero.Error) {
	return w.user(ctx, func(ctx context.Context, key any) (*models.SavedUser, ero.Error) {
		return w.Storage.UserById(ctx, key.(int64))
	}, id, "id")
}

type getUserFn func(ctx context.Context, key any) (*models.SavedUser, ero.Error)

// Throws;
//
//	ErrInvalidCredentials
//	ErrUserDeleted
//	Internal
func (w *Wrapper) user(ctx context.Context, get getUserFn, key any, keyType string) (*models.SavedUser, ero.Error) {
	const op = "wrap.Wrapper.user"
	log := w.log.With(slog.String("op", op), slog.Any(keyType, key))

	saved, err := get(ctx, key)
	switch {
	case ero.Has(err, storage.ErrEmptyResult):
		log.Error("no user found")
		return nil, ero.NewClient(ErrInvalidCredentials)
	case err != nil:
		log.Error("getting error", slog.String("error", err.Error()))
		return nil, ero.InternalFrom(op, err)
	}

	if saved.IsDeleted() {
		log.Error("user deleted")
		return saved, ero.NewClient(ErrUserDeleted)
	}

	return saved, nil
}
