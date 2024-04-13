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

func (pg *Pg) SaveSession(ctx context.Context, session *models.Session) (*models.SavedSession, error) {
	const op = "postgres.Pg.SaveSession"

	stmt, err := pg.db.PreparexContext(ctx, `
		INSERT INTO sessions (signup_fk, ip, browser, os) VALUES ($1, $2, $3, $4)
		RETURNING session_uuid, signup_fk, ip, browser, os, at`)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, session.SignupId, session.Info.Ip, session.Info.Browser, session.Info.OS)

	err = row.Err()
	if pgxerr.Is(err, pgerrcode.UniqueViolation) {
		return nil, storage.ErrUniqueConstraint
	}
	if pgxerr.Is(err, pgerrcode.ForeignKeyViolation) {
		return nil, storage.ErrFKConstraint
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

func (pg *Pg) SessionByUuid(ctx context.Context, uuid string) (*models.SavedSession, error) {
	return pg.whereSession(ctx, "session_uuid = $1", uuid)
}

func (pg *Pg) SessionByInfo(ctx context.Context, signupId int64, info models.SessionInfo) (*models.SavedSession, error) {
	return pg.whereSession(ctx, "browser = $1 AND ip = $2 AND os = $3 AND signup_fk = $4", info.Browser, info.Ip, info.OS, signupId)
}

func (pg *Pg) whereSession(ctx context.Context, where string, args ...any) (*models.SavedSession, error) {
	const op = "postgres.Pg.Session"

	s := fmt.Sprintf(`
	SELECT *
	FROM sessions
	WHERE %s`, where)

	stmt, err := pg.db.PreparexContext(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}

	saved := &models.SavedSession{}
	err = stmt.GetContext(ctx, saved, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrEmptyResult
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return saved, nil
}

func (pg *Pg) TerminateSession(ctx context.Context, uuid string) error {
	return pg.updateSession(ctx, "terminated_at = NOW()", "session_uuid = $1", uuid)
}

func (pg *Pg) ReviveSession(ctx context.Context, uuid string) error {
	return pg.updateSession(ctx, "terminated_at = NULL", "session_uuid = $1", uuid)
}

func (pg *Pg) updateSession(ctx context.Context, set, where string, args ...any) error {
	const op = "postgres.Pg.TerminateSession"

	s := fmt.Sprintf(`
	UPDATE sessions
	SET %s
	WHERE %s`, set, where)

	stmt, err := pg.db.PreparexContext(ctx, s)
	if err != nil {
		return fmt.Errorf("preparex %s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, args...)
	err = row.Err()
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrEmptyResult
	}
	if pgxerr.Is(err, pgerrcode.UniqueViolation) {
		return storage.ErrUniqueConstraint
	}

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (pg *Pg) DeleteSession(ctx context.Context, uuid string) error {
	const op = "postgres.Pg.DeleteSession"

	stmt, err := pg.db.PreparexContext(ctx, `
	DELETE FROM sessions
	WHERE session_uuid = $1`)
	if err != nil {
		return fmt.Errorf("preparex %s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, uuid)

	err = row.Err()
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrEmptyResult
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
