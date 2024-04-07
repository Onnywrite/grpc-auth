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

func (pg *Pg) UserBy(ctx context.Context, identifier models.UserIdentifier) (*models.SavedUser, error) {
	const op = "postgres.Pg.userBy"

	s := fmt.Sprintf(`
	SELECT user_id, login, email, phone, password, deleted_at
	FROM users
	WHERE %s = $1
	`, identifier.Key)
	
	stmt, err := pg.db.PreparexContext(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}
	
	u := &models.SavedUser{}
	err = stmt.GetContext(ctx, u, identifier.Value)
	
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrEmptyResult
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	return u, nil
}
