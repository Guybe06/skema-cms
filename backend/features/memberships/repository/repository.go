package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"skema-api/features/memberships/types"
)

type Repository struct{ db *pgxpool.Pool }

func New(db *pgxpool.Pool) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, m *types.Membership) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO memberships
		 (id, organization_id, user_id, email, role, status, invited_by, token_hash, expires_at, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		m.ID, m.OrganizationID, uuidParam(m.UserID), m.Email,
		m.Role, m.Status, uuidParam(m.InvitedBy), textParam(m.TokenHash),
		m.ExpiresAt, m.CreatedAt, m.UpdatedAt,
	)
	return err
}

func (r *Repository) FindByOrgAndUser(ctx context.Context, orgID, userID string) (*types.Membership, error) {
	return r.scanOne(ctx,
		`SELECT id, organization_id, COALESCE(user_id::text,''), email, role, status,
		        COALESCE(invited_by::text,''), COALESCE(token_hash,''), expires_at, created_at, updated_at
		 FROM memberships WHERE organization_id=$1 AND user_id=$2`, orgID, userID)
}

func (r *Repository) FindByOrgAndEmail(ctx context.Context, orgID, email string) (*types.Membership, error) {
	return r.scanOne(ctx,
		`SELECT id, organization_id, COALESCE(user_id::text,''), email, role, status,
		        COALESCE(invited_by::text,''), COALESCE(token_hash,''), expires_at, created_at, updated_at
		 FROM memberships WHERE organization_id=$1 AND email=$2`, orgID, email)
}

func (r *Repository) FindByToken(ctx context.Context, tokenHash string) (*types.Membership, error) {
	return r.scanOne(ctx,
		`SELECT id, organization_id, COALESCE(user_id::text,''), email, role, status,
		        COALESCE(invited_by::text,''), COALESCE(token_hash,''), expires_at, created_at, updated_at
		 FROM memberships WHERE token_hash=$1`, tokenHash)
}

func (r *Repository) ListByOrg(ctx context.Context, orgID string) ([]*types.Membership, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, organization_id, COALESCE(user_id::text,''), email, role, status,
		        COALESCE(invited_by::text,''), COALESCE(token_hash,''), expires_at, created_at, updated_at
		 FROM memberships WHERE organization_id=$1 ORDER BY created_at ASC`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*types.Membership
	for rows.Next() {
		m := &types.Membership{}
		if err := rows.Scan(&m.ID, &m.OrganizationID, &m.UserID, &m.Email,
			&m.Role, &m.Status, &m.InvitedBy, &m.TokenHash, &m.ExpiresAt,
			&m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

func (r *Repository) UpdateRole(ctx context.Context, id, role string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE memberships SET role=$1, updated_at=$2 WHERE id=$3`,
		role, time.Now(), id)
	return err
}

func (r *Repository) AcceptInvite(ctx context.Context, id, userID string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE memberships SET user_id=$1, status='active', token_hash=NULL, expires_at=NULL, updated_at=$2 WHERE id=$3`,
		userID, time.Now(), id)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM memberships WHERE id=$1`, id)
	return err
}

func (r *Repository) scanOne(ctx context.Context, query string, args ...any) (*types.Membership, error) {
	m := &types.Membership{}
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&m.ID, &m.OrganizationID, &m.UserID, &m.Email,
		&m.Role, &m.Status, &m.InvitedBy, &m.TokenHash, &m.ExpiresAt,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return m, err
}
