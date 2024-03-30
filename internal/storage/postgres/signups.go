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

	row := pg.db.QueryRowxContext(
		ctx,
		fmt.Sprintf(`INSERT INTO signups (user_fk, service_fk, at) VALUES (%d, %d, NOW())`,
			signup.UserId, signup.ServiceId),
	)
	err := row.Err()
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
