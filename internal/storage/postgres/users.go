package postgres

import (
	"context"
	"fmt"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/models"
)

func (pg *Pg) SaveUser(ctx context.Context, user *models.User) (*models.SavedUser, ero.Error) {
	const op = "postgres.Pg.SaveUser"

	stmt, err := pg.db.PreparexContext(ctx, `
		INSERT INTO users (login, email, phone, password)
		VALUES ($1, $2, $3, $4)
		RETURNING user_id, login, email, phone, password, deleted_at`)
	if err != nil {
		return nil, preparingError(err, op)
	}

	row := stmt.QueryRowxContext(ctx, user.Login, user.Email, user.Phone, user.Password)
	if err = row.Err(); err != nil {
		return nil, queryError(err, op)
	}

	u := &models.SavedUser{}
	err = row.StructScan(u)
	if err != nil {
		return nil, scanningError(err, op)
	}

	return u, nil
}

func (pg *Pg) UserById(ctx context.Context, id int64) (*models.SavedUser, ero.Error) {
	return pg.whereUser(ctx, "user_id = $1", id)
}

func (pg *Pg) UserByLogin(ctx context.Context, login string) (*models.SavedUser, ero.Error) {
	return pg.whereUser(ctx, "login = $1", login)
}

func (pg *Pg) UserByEmail(ctx context.Context, email string) (*models.SavedUser, ero.Error) {
	return pg.whereUser(ctx, "email = $1", email)
}

func (pg *Pg) UserByPhone(ctx context.Context, phone string) (*models.SavedUser, ero.Error) {
	return pg.whereUser(ctx, "phone = $1", phone)
}

func (pg *Pg) whereUser(ctx context.Context, where string, args ...any) (*models.SavedUser, ero.Error) {
	const op = "postgres.Pg.userBy"

	stmt, err := pg.db.PreparexContext(ctx, fmt.Sprintf(`
		SELECT user_id, login, email, phone, deleted_at
		FROM users
		WHERE %s`, where))
	if err != nil {
		return nil, preparingError(err, op)
	}

	u := &models.SavedUser{}
	err = stmt.GetContext(ctx, u, args...)
	if err != nil {
		return nil, queryError(err, op)
	}

	return u, nil
}
