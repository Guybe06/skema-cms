package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) IsMember(ctx context.Context, orgID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM memberships WHERE organization_id=$1 AND user_id=$2 AND status='active')`,
		orgID, userID).Scan(&exists)
	return exists, err
}

func (r *Repository) IsAdminOrOwner(ctx context.Context, orgID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(
		   SELECT 1 FROM memberships
		   WHERE organization_id=$1 AND user_id=$2 AND status='active' AND role='admin'
		 ) OR EXISTS(
		   SELECT 1 FROM organizations WHERE id=$1 AND owner_id=$2
		 )`, orgID, userID).Scan(&exists)
	return exists, err
}

func uuidParam(s string) pgtype.UUID {
	var u pgtype.UUID
	if s == "" {
		return u
	}
	_ = u.Scan(s)
	return u
}

func textParam(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}
