package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"skema-api/features/accounts/types"
)

func (r *Repository) CreateSession(ctx context.Context, s *types.Session) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO sessions (id, user_id, token_hash, expires_at, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		s.ID, s.UserID, s.TokenHash, s.ExpiresAt, s.CreatedAt,
	)
	return err
}

func (r *Repository) FindSessionByTokenHash(ctx context.Context, tokenHash string) (*types.Session, error) {
	s := &types.Session{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, token_hash, expires_at, created_at
		 FROM sessions WHERE token_hash = $1 AND expires_at > $2`,
		tokenHash, time.Now(),
	).Scan(&s.ID, &s.UserID, &s.TokenHash, &s.ExpiresAt, &s.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return s, err
}

func (r *Repository) DeleteSession(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, id)
	return err
}

func (r *Repository) DeleteUserSessions(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}
