package postgres

import (
	"context"
	"fmt"

	"github.com/Onnywrite/grpc-auth/internal/lib/pgxerr"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/Onnywrite/grpc-auth/internal/storage"
	"github.com/jackc/pgerrcode"
)

func (pg *Pg) SaveUser(ctx context.Context, user *models.User) (*models.SavedUser, error) {
	const op = "postgres.Pg.SaveUser"

	result, err := pg.db.ExecContext(ctx, fmt.Sprintf(
		`INSERT INTO users (login, email, phone, password) VALUES ('%s','%s','%s','%s');`,
		user.Login, user.Email, user.Phone, user.Password))

	if pgxerr.Is(err, pgerrcode.UniqueViolation) {
		return nil, storage.ErrUserExists
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return user.Saved(id), nil
}
