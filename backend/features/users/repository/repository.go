package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"skema-api/features/users/types"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByID(ctx context.Context, id string) (*types.User, error) {
	u := &types.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, first_name, last_name, email_verified, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (r *Repository) UpdateProfile(ctx context.Context, id, firstName, lastName string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET first_name = $1, last_name = $2, updated_at = $3 WHERE id = $4`,
		firstName, lastName, time.Now(), id,
	)
	return err
}

func (r *Repository) FindPasswordHash(ctx context.Context, id string) (string, error) {
	var hash string
	err := r.db.QueryRow(ctx, `SELECT password_hash FROM users WHERE id = $1`, id).Scan(&hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	return hash, err
}

func (r *Repository) UpdatePassword(ctx context.Context, id, hash string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`,
		hash, time.Now(), id,
	)
	return err
}
