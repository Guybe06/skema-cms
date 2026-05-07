package repository

import "context"

func (r *Repository) IsMember(ctx context.Context, orgID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM memberships WHERE organization_id = $1 AND user_id = $2 AND status = 'active')`,
		orgID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *Repository) FindBySlugID(ctx context.Context, slug string) (string, error) {
	org, err := r.FindBySlug(ctx, slug)
	if err != nil || org == nil {
		return "", err
	}
	return org.ID, nil
}

func (r *Repository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM organizations WHERE slug = $1)`, slug,
	).Scan(&exists)
	return exists, err
}
