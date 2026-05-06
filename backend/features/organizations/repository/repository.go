package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"skema-api/features/organizations/types"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, org *types.Organization) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO organizations (id, name, slug, owner_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		org.ID, org.Name, org.Slug, org.OwnerID, org.CreatedAt, org.UpdatedAt,
	)
	return err
}

func (r *Repository) FindBySlug(ctx context.Context, slug string) (*types.Organization, error) {
	org := &types.Organization{}
	err := r.db.QueryRow(ctx,
		`SELECT id, name, slug, owner_id, created_at, updated_at
		 FROM organizations WHERE slug = $1`, slug,
	).Scan(&org.ID, &org.Name, &org.Slug, &org.OwnerID, &org.CreatedAt, &org.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return org, err
}

func (r *Repository) FindByID(ctx context.Context, id string) (*types.Organization, error) {
	org := &types.Organization{}
	err := r.db.QueryRow(ctx,
		`SELECT id, name, slug, owner_id, created_at, updated_at
		 FROM organizations WHERE id = $1`, id,
	).Scan(&org.ID, &org.Name, &org.Slug, &org.OwnerID, &org.CreatedAt, &org.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return org, err
}

func (r *Repository) FindByOwner(ctx context.Context, ownerID string) ([]*types.Organization, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, slug, owner_id, created_at, updated_at
		 FROM organizations WHERE owner_id = $1 ORDER BY created_at DESC`, ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []*types.Organization
	for rows.Next() {
		org := &types.Organization{}
		if err := rows.Scan(&org.ID, &org.Name, &org.Slug, &org.OwnerID, &org.CreatedAt, &org.UpdatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}
	return orgs, rows.Err()
}

func (r *Repository) Update(ctx context.Context, id, name, slug string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE organizations SET name = $1, slug = $2, updated_at = $3 WHERE id = $4`,
		name, slug, time.Now(), id,
	)
	return err
}

func (r *Repository) TransferOwnership(ctx context.Context, id, newOwnerID string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE organizations SET owner_id = $1, updated_at = $2 WHERE id = $3`,
		newOwnerID, time.Now(), id,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM organizations WHERE id = $1`, id)
	return err
}

func (r *Repository) IsMember(ctx context.Context, orgID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM memberships WHERE organization_id = $1 AND user_id = $2 AND status = 'active')`,
		orgID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM organizations WHERE slug = $1)`, slug,
	).Scan(&exists)
	return exists, err
}
