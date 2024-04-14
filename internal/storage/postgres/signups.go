package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Onnywrite/grpc-auth/internal/lib/pgxerr"
	"github.com/Onnywrite/grpc-auth/internal/models"
	storage "github.com/Onnywrite/grpc-auth/internal/storage/common"
	"github.com/jackc/pgerrcode"
)

func (pg *Pg) SaveSignup(ctx context.Context, signup models.Signup) (*models.SavedSignup, error) {
	const op = "postgres.Pg.SaveSignup"

	stmt, err := pg.db.PreparexContext(ctx, `
		INSERT INTO signups (user_fk, service_fk)
		VALUES ($1, $2)
		RETURNING signup_id, user_fk, service_fk, at, banned_at`)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, signup.UserId, signup.ServiceId)
	err = row.Err()
	if pgxerr.Is(err, pgerrcode.UniqueViolation) {
		return nil, storage.ErrUniqueConstraint
	}
	if pgxerr.Is(err, pgerrcode.ForeignKeyViolation) {
		return nil, storage.ErrFKConstraint
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	su := &models.SavedSignup{}
	row.StructScan(su)

	return su, nil
}

func (pg *Pg) SignupById(ctx context.Context, id int64) (*models.SavedSignup, error) {
	return pg.whereSignup(ctx, "signup_id = $1", id)
}

func (pg *Pg) SignupByServiceAndUser(ctx context.Context, serviceId, userId int64) (*models.SavedSignup, error) {
	return pg.whereSignup(ctx, "service_fk = $1 AND user_fk = $2", serviceId, userId)
}

func (pg *Pg) whereSignup(ctx context.Context, where string, args ...any) (*models.SavedSignup, error) {
	const op = "postgres.Pg.Signup"

	s := fmt.Sprintf(`
	SELECT signup_id, user_fk, service_fk, at, banned_at
	FROM signups
	WHERE %s`, where)

	stmt, err := pg.db.PreparexContext(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}

	su := &models.SavedSignup{}
	err = stmt.GetContext(ctx, su, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrEmptyResult
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return su, nil
}
