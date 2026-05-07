package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"skema-api/features/apikeys/types"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, k *types.APIKey) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO api_keys
		 (id, organization_id, name, key_hash, key_prefix, permissions, allowed_collections, expires_at, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		k.ID, k.OrganizationID, k.Name, k.KeyHash, k.KeyPrefix,
		k.Permissions, nullableJSON(k.AllowedCollections), k.ExpiresAt, k.CreatedAt,
	)
	return err
}

func (r *Repository) FindByHash(ctx context.Context, hash string) (*types.APIKey, error) {
	return r.scanOne(ctx,
		`SELECT id, organization_id, name, key_hash, key_prefix, permissions,
		        allowed_collections, expires_at, last_used_at, created_at
		 FROM api_keys WHERE key_hash=$1`, hash)
}

func (r *Repository) FindByID(ctx context.Context, id string) (*types.APIKey, error) {
	return r.scanOne(ctx,
		`SELECT id, organization_id, name, key_hash, key_prefix, permissions,
		        allowed_collections, expires_at, last_used_at, created_at
		 FROM api_keys WHERE id=$1`, id)
}

func (r *Repository) ListByOrg(ctx context.Context, orgID string) ([]*types.APIKey, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, organization_id, name, key_hash, key_prefix, permissions,
		        allowed_collections, expires_at, last_used_at, created_at
		 FROM api_keys WHERE organization_id=$1 ORDER BY created_at DESC`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*types.APIKey
	for rows.Next() {
		k := &types.APIKey{}
		if err := scanFields(rows.Scan, k); err != nil {
			return nil, err
		}
		list = append(list, k)
	}
	return list, rows.Err()
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM api_keys WHERE id=$1`, id)
	return err
}

func (r *Repository) TouchLastUsed(ctx context.Context, id string) {
	r.db.Exec(context.Background(), `UPDATE api_keys SET last_used_at=$1 WHERE id=$2`, time.Now(), id)
}

func (r *Repository) scanOne(ctx context.Context, query string, args ...any) (*types.APIKey, error) {
	k := &types.APIKey{}
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&k.ID, &k.OrganizationID, &k.Name, &k.KeyHash, &k.KeyPrefix,
		&k.Permissions, &k.AllowedCollections, &k.ExpiresAt, &k.LastUsedAt, &k.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return k, err
}

func scanFields(scan func(...any) error, k *types.APIKey) error {
	return scan(
		&k.ID, &k.OrganizationID, &k.Name, &k.KeyHash, &k.KeyPrefix,
		&k.Permissions, &k.AllowedCollections, &k.ExpiresAt, &k.LastUsedAt, &k.CreatedAt,
	)
}

func nullableJSON(raw []byte) any {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	return raw
}
