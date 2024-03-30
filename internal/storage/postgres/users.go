package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Onnywrite/grpc-auth/internal/lib/pgxerr"
	"github.com/Onnywrite/grpc-auth/internal/lib/sqlf"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/Onnywrite/grpc-auth/internal/storage"
	"github.com/jackc/pgerrcode"
)

func (pg *Pg) SaveUser(ctx context.Context, user *models.User) (*models.SavedUser, error) {
	const op = "postgres.Pg.SaveUser"

	row := pg.db.QueryRowxContext(ctx,
		sqlf.SQLFormat(`INSERT INTO users (login, email, phone, password)
		VALUES (%s, %s, %s, %s)
		RETURNING user_id, login, email, phone, password`,
			user.Login, user.Email, user.Phone, user.Password))
	err := row.Err()
	if pgxerr.Is(err, pgerrcode.UniqueViolation) {
		return nil, storage.ErrUserExists
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	u := &models.SavedUser{}
	err = row.StructScan(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (pg *Pg) UserById(ctx context.Context, id int) (u *models.SavedUser, err error) {
	return pg.userBy(ctx, "user_id", id)
}

func (pg *Pg) UserByLogin(ctx context.Context, login string) (u *models.SavedUser, err error) {
	return pg.userBy(ctx, "login", login)
}

func (pg *Pg) UserByEmail(ctx context.Context, email string) (u *models.SavedUser, err error) {
	return pg.userBy(ctx, "email", email)
}

func (pg *Pg) UserByPhone(ctx context.Context, phone string) (u *models.SavedUser, err error) {
	return pg.userBy(ctx, "phone", phone)
}

func (pg *Pg) userBy(ctx context.Context, prop string, val any) (*models.SavedUser, error) {
	const op = "postgres.Pg.userBy"

	u := &models.SavedUser{}
	err := pg.db.GetContext(ctx, u,
		sqlf.SQLFormat(fmt.Sprintf(`SELECT user_id, login, email, phone, password
			FROM users WHERE %s`, prop)+`= %s`, val),
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}
