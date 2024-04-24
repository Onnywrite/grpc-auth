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
	sqlerrToErr = map[string]*ero.StorageError{
		pgerrcode.UniqueViolation:     storage.ErrUniqueConstraint,
		pgerrcode.ForeignKeyViolation: storage.ErrFKConstraint,
		sql.ErrNoRows.Error():         storage.ErrEmptyResult,
	}
)

func pgerr(anyerr error, op string) error {
	pgErr := &pgconn.PgError{}
	var strerr string
	if errors.As(anyerr, &pgErr) {
		strerr = pgErr.Code
	} else {
		strerr = anyerr.Error()
	}

	if err, ok := sqlerrToErr[strerr]; ok {
		err.WithMethod(op)
		return err
	}

	return ero.NewStorage(anyerr.Error()).WithMethod(op)
}

func preperr(anyerr error, op string) error {
	return ero.NewStorage(anyerr.Error()).WithMethod(op).SetCode(ero.CodeStorageInternal)
}

func scanerr(anyerr error, op string) error {
	return ero.NewStorage(anyerr.Error()).WithMethod(op).SetCode(ero.CodeStorageInternal)
}
