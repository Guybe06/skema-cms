package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"skema-api/features/collections/types"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, c *types.Collection) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO collections
		 (id, connection_id, organization_id, name, table_name, display_name, description, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		c.ID, c.ConnectionID, c.OrganizationID, c.Name, c.TableName,
		c.DisplayName, c.Description, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *Repository) FindByID(ctx context.Context, id string) (*types.Collection, error) {
	c := &types.Collection{}
	err := r.db.QueryRow(ctx,
		`SELECT id, connection_id, organization_id, name, table_name,
		        COALESCE(display_name,''), COALESCE(description,''), created_at, updated_at
		 FROM collections WHERE id=$1`, id,
	).Scan(&c.ID, &c.ConnectionID, &c.OrganizationID, &c.Name, &c.TableName,
		&c.DisplayName, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return c, err
}

func (r *Repository) FindByOrgAndTable(ctx context.Context, orgID, tableName string) (*types.Collection, error) {
	c := &types.Collection{}
	err := r.db.QueryRow(ctx,
		`SELECT id, connection_id, organization_id, name, table_name,
		        COALESCE(display_name,''), COALESCE(description,''), created_at, updated_at
		 FROM collections WHERE organization_id=$1 AND table_name=$2`, orgID, tableName,
	).Scan(&c.ID, &c.ConnectionID, &c.OrganizationID, &c.Name, &c.TableName,
		&c.DisplayName, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return c, err
}

func (r *Repository) ListByOrg(ctx context.Context, orgID string) ([]*types.Collection, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, connection_id, organization_id, name, table_name,
		        COALESCE(display_name,''), COALESCE(description,''), created_at, updated_at
		 FROM collections WHERE organization_id=$1 ORDER BY created_at ASC`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*types.Collection
	for rows.Next() {
		c := &types.Collection{}
		if err := rows.Scan(&c.ID, &c.ConnectionID, &c.OrganizationID, &c.Name, &c.TableName,
			&c.DisplayName, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *Repository) Update(ctx context.Context, id, name, displayName, description string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE collections SET name=$1, display_name=$2, description=$3, updated_at=$4 WHERE id=$5`,
		name, displayName, description, time.Now(), id,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM collections WHERE id=$1`, id)
	return err
}

func (r *Repository) TableExists(ctx context.Context, connectionID, tableName string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM collections WHERE connection_id=$1 AND table_name=$2)`,
		connectionID, tableName,
	).Scan(&exists)
	return exists, err
}
