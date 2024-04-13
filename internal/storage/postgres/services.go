package postgres

import (
	"context"
	"fmt"

	"github.com/Onnywrite/grpc-auth/internal/lib/pgxerr"
	"github.com/Onnywrite/grpc-auth/internal/models"
	storage "github.com/Onnywrite/grpc-auth/internal/storage/common"
	"github.com/jackc/pgerrcode"
)

func (pg *Pg) SaveService(ctx context.Context, service *models.Service) (*models.SavedService, error) {
	const op = "postgres.Pg.SaveService"

	stmt, err := pg.db.PreparexContext(ctx, `
		INSERT INTO services (name, owner_fk)
		VALUES $1, $2
		RETURNING service_id, owner_fk, name`)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}

	row := stmt.QueryRowxContext(ctx, service.Name, service.OwnerId)
	err = row.Err()
	if pgxerr.Is(err, pgerrcode.UniqueViolation) {
		return nil, storage.ErrUniqueConstraint
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	u := &models.SavedService{}
	err = row.StructScan(u)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}

func (pg *Pg) ServiceById(ctx context.Context, id int64) (*models.SavedService, error) {
	return pg.whereService(ctx, "service_id = $1", id)
}

func (pg *Pg) whereService(ctx context.Context, where string, args ...any) (*models.SavedService, error) {
	const op = "postgres.Pg.serviceBy"

	s := fmt.Sprintf(`
		SELECT *
		FROM services
		WHERE %s`, where)

	stmt, err := pg.db.PreparexContext(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("preparex %s: %w", op, err)
	}

	saved := &models.SavedService{}
	err = stmt.GetContext(ctx, saved, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return saved, nil
}
