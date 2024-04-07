package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Onnywrite/grpc-auth/internal/lib/pgxerr"
	"github.com/Onnywrite/grpc-auth/internal/models"
	storage "github.com/Onnywrite/grpc-auth/internal/storage/common"
	"github.com/jackc/pgerrcode"
)

func (pg *Pg) SaveSession(ctx context.Context, session *models.Session) (*models.SavedSession, error) {
	const op = "postgres.Pg.SaveSession"

	stmt, err := pg.db.PreparexContext(ctx,
		`INSERT INTO sessions (signup_fk, ip, browser, os) VALUES (
		(
			SELECT signup_id
			FROM signups
			WHERE user_fk = $1 AND service_fk = $2
		), $3, $4, $5)
		RETURNING session_uuid, signup_fk, ip, browser, os, at`)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, session.UserId, session.ServiceId, session.Info.Ip, session.Info.Browser, session.Info.OS)

	err = row.Err()
	if pgxerr.Is(err, pgerrcode.UniqueViolation) {
		return nil, storage.ErrUniqueConstraint
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ss := &models.SavedSession{}
	err = row.StructScan(ss)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ss, nil
}

func (pg *Pg) DeleteSession(ctx context.Context, uuid string) error {
	const op = "postgres.Pg.DeleteSession"

	s, args, err := sq.Delete("sessions").
		Where(sq.Eq{"session_uuid": uuid}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("squirrel %s: %w", op, err)
	}

	stmt, err := pg.db.PreparexContext(ctx, s)
	if err != nil {
		return fmt.Errorf("preparex %s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, args...)

	err = row.Err()
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrEmptyResult
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (pg *Pg) SessionById(ctx context.Context, uuid string) (*models.SavedSession, error) {
	const op = "postgres.Pg.SessionById"

	s, args, err := sq.Select("sessions").
		Where(sq.Eq{"session_uuid": uuid}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("squirrel %s: %w", op, err)
	}

	stmt, err := pg.db.PreparexContext(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}

	session := &models.SavedSession{}
	err = stmt.GetContext(ctx, session, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrEmptyResult
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}

func (pg *Pg) Session(ctx context.Context, session *models.Session) (*models.SavedSession, error) {
	const op = "postgres.Pg.Session"

	stmt, err := pg.db.PreparexContext(ctx, `
	SELECT *
	FROM sessions
	WHERE signup_fk IN(
		SELECT signup_id
		FROM signups
		WHERE user_fk = $1 AND service_fk = $2
	) AND browser = $3 AND ip = $4 AND os = $5
	`)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}

	saved := &models.SavedSession{}
	err = stmt.GetContext(ctx, saved, session.UserId, session.ServiceId, session.Info.Browser, session.Info.Ip, session.Info.OS)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrEmptyResult
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return saved, nil
}

func (pg *Pg) TerminateSession(ctx context.Context, uuid string) error {
	const op = "postgres.Pg.TerminateSession"

	s, args, err := sq.Update("sessions").
		Set("terminated_at", time.Now()).
		Where(sq.And{sq.Eq{"session_uuid": uuid}, sq.Eq{"terminated_at": nil}}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("squirrel %s: %w", op, err)
	}

	stmt, err := pg.db.PreparexContext(ctx, s)
	if err != nil {
		return fmt.Errorf("preparex %s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, args...)
	err = row.Err()
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrEmptyResult
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
