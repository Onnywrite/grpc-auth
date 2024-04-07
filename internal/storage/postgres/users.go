package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Onnywrite/grpc-auth/internal/lib/pgxerr"
	"github.com/Onnywrite/grpc-auth/internal/models"
	storage "github.com/Onnywrite/grpc-auth/internal/storage/common"
	"github.com/jackc/pgerrcode"
)

func (pg *Pg) SaveUser(ctx context.Context, user *models.User) (*models.SavedUser, error) {
	const op = "postgres.Pg.SaveUser"

	s, args, err := sq.Insert("users").Columns("login", "email", "phone", "password").
		Values(user.Login, user.Email, user.Phone, user.Password).
		Suffix(`RETURNING user_id, login, email, phone, password`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("squirrel %s: %w", op, err)
	}

	stmt, err := pg.db.PreparexContext(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, args...)
	err = row.Err()
	if pgxerr.Is(err, pgerrcode.UniqueViolation) {
		return nil, storage.ErrUserExists
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

func (pg *Pg) UserById(ctx context.Context, id int64) (u *models.SavedUser, err error) {
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

	s, args, err := sq.Select("user_id", "login", "email", "phone", "password").
		From("users").
		Where(sq.Eq{prop: val}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("squirrel %s: %w", op, err)
	}

	stmt, err := pg.db.PreparexContext(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}

	u := &models.SavedUser{}
	err = stmt.GetContext(ctx, u, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}
