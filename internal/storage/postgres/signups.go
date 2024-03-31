package postgres

import (
	"context"
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
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if pgxerr.Is(err, pgerrcode.UniqueViolation) {
		return storage.ErrSignupExists
	}
	if pgxerr.Is(err, pgerrcode.ForeignKeyViolation) {
		return storage.ErrNoSuchPrimaryKey
	}

	return nil
}
