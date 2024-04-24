package postgres

import (
	"context"
	"fmt"

	"github.com/Onnywrite/grpc-auth/internal/models"
)

func (pg *Pg) SaveService(ctx context.Context, service *models.Service) (*models.SavedService, error) {
	const op = "postgres.Pg.SaveService"

	stmt, err := pg.db.PreparexContext(ctx, `
		INSERT INTO services (name, owner_fk)
		VALUES ($1, $2)
		RETURNING service_id, owner_fk, name, deleted_at`)
	if err != nil {
		return nil, preperr(err, op)
	}

	row := stmt.QueryRowxContext(ctx, service.Name, service.OwnerId)
	if err = row.Err(); err != nil {
		return nil, pgerr(err, op)
	}

	u := &models.SavedService{}
	err = row.StructScan(u)
	if err != nil {
		return nil, scanerr(err, op)
	}

	return u, nil
}

func (pg *Pg) ServiceById(ctx context.Context, id int64) (*models.SavedService, error) {
	return pg.whereService(ctx, "service_id = $1", id)
}

func (pg *Pg) whereService(ctx context.Context, where string, args ...any) (*models.SavedService, error) {
	const op = "postgres.Pg.serviceBy"

	stmt, err := pg.db.PreparexContext(ctx, fmt.Sprintf(`
		SELECT *
		FROM services
		WHERE %s`, where))
	if err != nil {
		return nil, preperr(err, op)
	}

	saved := &models.SavedService{}
	err = stmt.GetContext(ctx, saved, args...)
	if err != nil {
		return nil, pgerr(err, op)
	}

	return saved, nil
}
