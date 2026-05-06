package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"skema-api/features/auth/types"
)

func (r *Repository) CreateVerificationToken(ctx context.Context, t *types.VerificationToken) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO verification_tokens (id, user_id, token_hash, type, expires_at, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		t.ID, t.UserID, t.TokenHash, t.Type, t.ExpiresAt, t.CreatedAt,
	)
	return err
}

func (r *Repository) FindVerificationToken(ctx context.Context, tokenHash, tokenType string) (*types.VerificationToken, error) {
	t := &types.VerificationToken{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, token_hash, type, expires_at, created_at
		 FROM verification_tokens WHERE token_hash = $1 AND type = $2 AND expires_at > $3`,
		tokenHash, tokenType, time.Now(),
	).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.Type, &t.ExpiresAt, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return t, err
}

func (r *Repository) DeleteVerificationToken(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM verification_tokens WHERE id = $1`, id)
	return err
}

func (r *Repository) DeleteTokensByUser(ctx context.Context, userID, tokenType string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM verification_tokens WHERE user_id = $1 AND type = $2`,
		userID, tokenType,
	)
	return err
}
