package postgres

import (
	"context"
	"fmt"

	"github.com/Onnywrite/grpc-auth/internal/lib/ero"
	"github.com/Onnywrite/grpc-auth/internal/models"
)

func (pg *Pg) SaveSession(ctx context.Context, session *models.Session) (*models.SavedSession, ero.Error) {
	const op = "postgres.Pg.SaveSession"

	stmt, err := pg.db.PreparexContext(ctx, `
		INSERT INTO sessions (service_fk, user_fk, ip, browser, os)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING session_uuid, service_fk, user_fk, ip, browser, os, at`)
	if err != nil {
		return nil, preparingError(err, op)
	}

	row := stmt.QueryRowxContext(ctx, session.ServiceId, session.UserId, session.Info.Ip, session.Info.Browser, session.Info.OS)
	if err = row.Err(); err != nil {
		return nil, queryError(err, op)
	}

	ss := &models.SavedSession{}
	err = row.StructScan(ss)
	if err != nil {
		return nil, scanningError(err, op)
	}

	return ss, nil
}

func (pg *Pg) SessionByUuid(ctx context.Context, uuid string) (*models.SavedSession, ero.Error) {
	return pg.whereSession(ctx, "session_uuid = $1", uuid)
}

func (pg *Pg) SessionByInfo(ctx context.Context, serviceId, userId int64, info models.SessionInfo) (*models.SavedSession, ero.Error) {
	ifNil := func(s *string) string {
		if s == nil {
			return "IS NULL"
		}
		return `= '` + *s + `'`
	}
	return pg.whereSession(ctx, `browser `+ifNil(info.Browser)+
		` AND ip `+ifNil(info.Ip)+
		` AND os `+ifNil(info.OS)+
		` AND service_fk = $1 AND user_fk = $2`,
		serviceId, userId)
}

func (pg *Pg) whereSession(ctx context.Context, where string, args ...any) (*models.SavedSession, ero.Error) {
	const op = "postgres.Pg.Session"

	stmt, err := pg.db.PreparexContext(ctx, fmt.Sprintf(`
		SELECT *
		FROM sessions
		WHERE %s`, where))
	if err != nil {
		return nil, preparingError(err, op)
	}

	saved := &models.SavedSession{}
	err = stmt.GetContext(ctx, saved, args...)
	if err != nil {
		return nil, queryError(err, op)
	}

	return saved, nil
}

func (pg *Pg) DeleteSession(ctx context.Context, uuid string) ero.Error {
	const op = "postgres.Pg.DeleteSession"

	stmt, err := pg.db.PreparexContext(ctx, `
		DELETE FROM sessions
		WHERE session_uuid = $1`)
	if err != nil {
		return preparingError(err, op)
	}

	row := stmt.QueryRowxContext(ctx, uuid)
	if err = row.Err(); err != nil {
		return queryError(err, op)
	}

	return nil
}
