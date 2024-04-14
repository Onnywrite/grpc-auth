package storage

import "github.com/Onnywrite/grpc-auth/internal/storage/postgres"

type Storage struct {
	*postgres.Pg
}

func New(conn string) (*Storage, error) {
	pg, err := postgres.NewPg(conn)
	if err != nil {
		return nil, err
	}

	return &Storage{
		pg,
	}, nil
}
