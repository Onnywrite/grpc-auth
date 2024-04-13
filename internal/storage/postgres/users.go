package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Onnywrite/grpc-auth/internal/lib/pgxerr"
	"github.com/Onnywrite/grpc-auth/internal/models"
	storage "github.com/Onnywrite/grpc-auth/internal/storage/common"
	"github.com/jackc/pgerrcode"
)

func (pg *Pg) SaveUser(ctx context.Context, user *models.User) (*models.SavedUser, error) {
	const op = "postgres.Pg.SaveUser"

	stmt, err := pg.db.PreparexContext(ctx, `
		INSERT INTO users (login, email, phone, password)
		VALUES $1, $2, $3, $4
		RETURNING user_id, login, email, phone, password`)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, user.Login, user.Email, user.Phone, user.Password)
	err = row.Err()
	if pgxerr.Is(err, pgerrcode.UniqueViolation) {
		return nil, storage.ErrUniqueConstraint
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	u := &models.SavedUser{}
	err = row.StructScan(u)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}

func (pg *Pg) UserById(ctx context.Context, id int64) (*models.SavedUser, error) {
	return pg.whereUser(ctx, "user_id = $1", id)
}

func (pg *Pg) UserByLogin(ctx context.Context, login string) (*models.SavedUser, error) {
	return pg.whereUser(ctx, "login = $1", login)
}

func (pg *Pg) UserByEmail(ctx context.Context, email string) (*models.SavedUser, error) {
	return pg.whereUser(ctx, "email = $1", email)
}

func (pg *Pg) UserByPhone(ctx context.Context, phone string) (*models.SavedUser, error) {
	return pg.whereUser(ctx, "phone = $1", phone)
}

func (pg *Pg) whereUser(ctx context.Context, where string, args ...any) (*models.SavedUser, error) {
	const op = "postgres.Pg.userBy"

	s := fmt.Sprintf(`
	SELECT user_id, login, email, phone, password, deleted_at
	FROM users
	WHERE %s`, where)

	stmt, err := pg.db.PreparexContext(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}

	u := &models.SavedUser{}
	err = stmt.GetContext(ctx, u, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrEmptyResult
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}
