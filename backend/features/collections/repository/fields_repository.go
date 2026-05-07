package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"skema-api/features/collections/types"
)

func (r *Repository) AddField(ctx context.Context, f *types.Field) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO collection_fields
		 (id, collection_id, name, column_name, type, required, is_unique, default_value, options, position, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		f.ID, f.CollectionID, f.Name, f.ColumnName, f.Type,
		f.Required, f.IsUnique, f.DefaultValue, nullableJSON(f.Options), f.Position, f.CreatedAt,
	)
	return err
}

func (r *Repository) ListFields(ctx context.Context, collectionID string) ([]*types.Field, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, collection_id, name, column_name, type, required, is_unique,
		        COALESCE(default_value,''), options, position, created_at
		 FROM collection_fields WHERE collection_id=$1 ORDER BY position ASC, created_at ASC`,
		collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []*types.Field
	for rows.Next() {
		f := &types.Field{}
		if err := rows.Scan(&f.ID, &f.CollectionID, &f.Name, &f.ColumnName, &f.Type,
			&f.Required, &f.IsUnique, &f.DefaultValue, &f.Options, &f.Position, &f.CreatedAt); err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return fields, rows.Err()
}

func (r *Repository) FindField(ctx context.Context, id string) (*types.Field, error) {
	f := &types.Field{}
	err := r.db.QueryRow(ctx,
		`SELECT id, collection_id, name, column_name, type, required, is_unique,
		        COALESCE(default_value,''), options, position, created_at
		 FROM collection_fields WHERE id=$1`, id,
	).Scan(&f.ID, &f.CollectionID, &f.Name, &f.ColumnName, &f.Type,
		&f.Required, &f.IsUnique, &f.DefaultValue, &f.Options, &f.Position, &f.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return f, err
}

func (r *Repository) ColumnExists(ctx context.Context, collectionID, columnName string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM collection_fields WHERE collection_id=$1 AND column_name=$2)`,
		collectionID, columnName,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) DeleteField(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM collection_fields WHERE id=$1`, id)
	return err
}

func nullableJSON(raw []byte) any {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	return raw
}
