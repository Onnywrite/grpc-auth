package postgres

import (
	"context"
	"fmt"

	"github.com/Onnywrite/grpc-auth/internal/lib/pgxerr"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/Onnywrite/grpc-auth/internal/storage"
	"github.com/jackc/pgerrcode"
)

func (pg *Pg) SaveSignup(ctx context.Context, signup models.Signup) error {
	const op = "postgres.Pg.SaveSignup"

	stmt, err := pg.db.PreparexContext(ctx, `INSERT INTO signups (user_fk, service_fk, at) VALUES ($1, $2)`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, signup.UserId, signup.ServiceId)
	err = row.Err()

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
