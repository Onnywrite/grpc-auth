package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Onnywrite/grpc-auth/internal/lib/pgxerr"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/Onnywrite/grpc-auth/internal/storage"
	"github.com/jackc/pgerrcode"
)

func (pg *Pg) SaveService(ctx context.Context, service *models.Service) (*models.SavedService, error) {
	const op = "postgres.Pg.SaveService"

	s, args, err := sq.Insert("services").
		Columns("name", "owner_fk").
		Values(service.Name, service.OwnerId).
		Suffix(`RETURNING service_id, owner_fk, name`).
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
		return nil, storage.ErrServiceExists
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

func (pg *Pg) Service(ctx context.Context, id int64) (*models.SavedService, error) {
	return pg.serviceBy(ctx, "service_id", id)
}

func (pg *Pg) serviceBy(ctx context.Context, prop string, val any) (*models.SavedService, error) {
	const op = "postgres.Pg.serviceBy"

	s, args, err := sq.Select("services").
		From("services").
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

	saved := &models.SavedService{}
	err = stmt.GetContext(ctx, saved, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return saved, nil
}
