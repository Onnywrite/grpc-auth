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
	"github.com/Onnywrite/grpc-auth/internal/storage"
	"github.com/google/uuid"
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
		RETURNING session_uuid, signup_fk, ip, browser, os, created_at`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, session.UserId, session.ServiceId, session.IP, session.Browser, session.OS)

	err = row.Err()
	if pgxerr.Is(err, pgerrcode.UniqueViolation) {
		return nil, storage.ErrSessionExists
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	saved := &models.SavedSession{}
	err = row.StructScan(saved)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return saved, nil
}

func (pg *Pg) DeleteSession(ctx context.Context, uuid uuid.UUID) error {
	const op = "postgres.Pg.DeleteSession"

	s, args, err := sq.Update("sessions").
		Set("terminated_at", time.Now()).
		Where(sq.Eq{"session_uuid": uuid.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := pg.db.PreparexContext(ctx, s)
	if err != nil {
		return fmt.Errorf("preparex %s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, args...)
	err = row.Err()
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrSessionNotFound
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
