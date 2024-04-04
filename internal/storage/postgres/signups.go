package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Onnywrite/grpc-auth/internal/lib/pgxerr"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/Onnywrite/grpc-auth/internal/storage"
	"github.com/jackc/pgerrcode"
)

func (pg *Pg) SaveSignup(ctx context.Context, signup models.Signup) error {
	const op = "postgres.Pg.SaveSignup"

	_, err := sq.Insert("signups").
		Columns("user_fk", "service_fk").
		Values(signup.UserId, signup.ServiceId).
		PlaceholderFormat(sq.Dollar).
		RunWith(pg.db).
		ExecContext(ctx)

	if pgxerr.Is(err, pgerrcode.UniqueViolation) {
		return storage.ErrSignupExists
	}
	if pgxerr.Is(err, pgerrcode.ForeignKeyViolation) {
		return storage.ErrNoSuchPrimaryKey
	}

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (pg *Pg) Signup(ctx context.Context, userId, serviceId int64) (*models.SavedSignup, error) {
	const op = "postgres.Pg.Signup"

	s, args, err := sq.Select("signup_id", "user_fk", "service_fk", "at", "banned_at").
	From("signups").
	Where(sq.Eq{"user_fk":userId, "service_fk":serviceId}).
	PlaceholderFormat(sq.Dollar).
	ToSql()

	if err != nil {
		return nil, fmt.Errorf("squirrel %s: %w", op, err)
	}

	stmt, err := pg.db.PreparexContext(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}

	su := &models.SavedSignup{}
	err = stmt.GetContext(ctx, su, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return su, nil
}