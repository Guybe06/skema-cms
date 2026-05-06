package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"skema-api/features/connections/types"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, rec *types.EncryptedRecord) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO connections
		 (id, organization_id, name, driver, host_encrypted, port_encrypted, database_encrypted, user_encrypted, password_encrypted, ssl_mode, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		rec.ID, rec.OrganizationID, rec.Name, rec.Driver,
		rec.HostEncrypted, rec.PortEncrypted, rec.DatabaseEncrypted, rec.UserEncrypted,
		rec.PasswordEncrypted, rec.SSLMode, rec.CreatedAt, rec.UpdatedAt,
	)
	return err
}

func (r *Repository) FindByID(ctx context.Context, id string) (*types.EncryptedRecord, error) {
	return r.scanOne(ctx,
		`SELECT id, organization_id, name, driver,
		        host_encrypted, port_encrypted, database_encrypted, user_encrypted,
		        password_encrypted, ssl_mode, created_at, updated_at
		 FROM connections WHERE id=$1`, id)
}

func (r *Repository) ListByOrg(ctx context.Context, orgID string) ([]*types.EncryptedRecord, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, organization_id, name, driver,
		        host_encrypted, port_encrypted, database_encrypted, user_encrypted,
		        password_encrypted, ssl_mode, created_at, updated_at
		 FROM connections WHERE organization_id=$1 ORDER BY created_at ASC`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*types.EncryptedRecord
	for rows.Next() {
		rec := &types.EncryptedRecord{}
		if err := rows.Scan(&rec.ID, &rec.OrganizationID, &rec.Name, &rec.Driver,
			&rec.HostEncrypted, &rec.PortEncrypted, &rec.DatabaseEncrypted, &rec.UserEncrypted,
			&rec.PasswordEncrypted, &rec.SSLMode, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, rec)
	}
	return list, rows.Err()
}

func (r *Repository) Update(ctx context.Context, rec *types.EncryptedRecord) error {
	_, err := r.db.Exec(ctx,
		`UPDATE connections
		 SET name=$1, host_encrypted=$2, port_encrypted=$3, database_encrypted=$4,
		     user_encrypted=$5, password_encrypted=$6, ssl_mode=$7, updated_at=$8
		 WHERE id=$9`,
		rec.Name, rec.HostEncrypted, rec.PortEncrypted, rec.DatabaseEncrypted,
		rec.UserEncrypted, rec.PasswordEncrypted, rec.SSLMode, time.Now(), rec.ID,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM connections WHERE id=$1`, id)
	return err
}

func (r *Repository) scanOne(ctx context.Context, query string, args ...any) (*types.EncryptedRecord, error) {
	rec := &types.EncryptedRecord{}
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&rec.ID, &rec.OrganizationID, &rec.Name, &rec.Driver,
		&rec.HostEncrypted, &rec.PortEncrypted, &rec.DatabaseEncrypted, &rec.UserEncrypted,
		&rec.PasswordEncrypted, &rec.SSLMode, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return rec, err
}
