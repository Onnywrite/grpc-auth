package storage

import (
	"fmt"

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
