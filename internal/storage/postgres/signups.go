package postgres

import (
	"context"
	"fmt"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/models"
)

func (pg *Pg) SaveSignup(ctx context.Context, signup models.Signup) (*models.SavedSignup, ero.Error) {
	const op = "postgres.Pg.SaveSignup"

	stmt, err := pg.db.PreparexContext(ctx, `
		INSERT INTO signups (user_fk, service_fk)
		VALUES ($1, $2)
		RETURNING user_fk, service_fk, at, banned_at`)
	if err != nil {
		return nil, preparingError(err, op)
	}

	row := stmt.QueryRowxContext(ctx, signup.UserId, signup.ServiceId)
	if err = row.Err(); err != nil {
		return nil, queryError(err, op)
	}

	su := &models.SavedSignup{}
	err = row.StructScan(su)
	if err != nil {
		return nil, scanningError(err, op)
	}

	return su, nil
}

func (pg *Pg) SignupByServiceAndUser(ctx context.Context, serviceId, userId int64) (*models.SavedSignup, ero.Error) {
	return pg.whereSignup(ctx, "service_fk = $1 AND user_fk = $2", serviceId, userId)
}

func (pg *Pg) whereSignup(ctx context.Context, where string, args ...any) (*models.SavedSignup, ero.Error) {
	const op = "postgres.Pg.Signup"

	stmt, err := pg.db.PreparexContext(ctx, fmt.Sprintf(`
		SELECT user_fk, service_fk, at, banned_at, deleted_at
		FROM signups
		WHERE %s`, where))
	if err != nil {
		return nil, preparingError(err, op)
	}

	su := &models.SavedSignup{}
	err = stmt.GetContext(ctx, su, args...)
	if err != nil {
		return nil, queryError(err, op)
	}

	return su, nil
}
