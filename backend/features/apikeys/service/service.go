package service

import (
	"context"

	"skema-api/features/apikeys/types"
	orgsrepo "skema-api/features/organizations/repository"
)

type repo interface {
	Create(ctx context.Context, k *types.APIKey) error
	FindByHash(ctx context.Context, hash string) (*types.APIKey, error)
	FindByID(ctx context.Context, id string) (*types.APIKey, error)
	ListByOrg(ctx context.Context, orgID string) ([]*types.APIKey, error)
	Delete(ctx context.Context, id string) error
	TouchLastUsed(ctx context.Context, id string)
}

type Service struct {
	repo     repo
	orgsRepo *orgsrepo.Repository
}

func New(repo repo, orgsRepo *orgsrepo.Repository) *Service {
	return &Service{repo: repo, orgsRepo: orgsRepo}
}

func (s *Service) FindByHash(ctx context.Context, hash string) (*types.APIKey, error) {
	return s.repo.FindByHash(ctx, hash)
}

func (s *Service) TouchLastUsed(ctx context.Context, id string) {
	s.repo.TouchLastUsed(ctx, id)
}
