package pgxerr

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func Is(err error, pgcode string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgcode
	}
	return false
}
