package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	storage "github.com/Onnywrite/grpc-auth/internal/storage/common"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Pg struct {
	db *sqlx.DB
}

func NewPg(conn string) (*Pg, error) {
	const op = "storage.NewPg"

	dbx, err := sqlx.Connect("pgx", conn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Pg{
		db: dbx,
	}, nil
}

func (pg *Pg) Disconnect() error {
	return pg.db.Close()
}

var (
	sqlerrToErr = map[string]string{
		pgerrcode.UniqueViolation:     storage.ErrUniqueConstraint,
		pgerrcode.ForeignKeyViolation: storage.ErrFKConstraint,
		sql.ErrNoRows.Error():         storage.ErrEmptyResult,
	}
)

func queryError(anyerr error, op string) ero.Error {
	pgErr := &pgconn.PgError{}
	var strerr string
	if errors.As(anyerr, &pgErr) {
		strerr = pgErr.Code
	} else {
		strerr = anyerr.Error()
	}

	if errStr, ok := sqlerrToErr[strerr]; ok {
		return ero.NewClient(ero.CodeUnknownClient, errStr)
	}

	return ero.NewServer(ero.CodeInternal, op, anyerr.Error())
}

func preparingError(anyerr error, op string) ero.Error {
	return ero.NewInternal(ero.CodeInternal, op, anyerr.Error())
}

func scanningError(anyerr error, op string) ero.Error {
	return ero.NewInternal(ero.CodeInternal, op, anyerr.Error())
}
