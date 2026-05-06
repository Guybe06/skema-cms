package service

import (
	"context"

	"skema-api/features/connections/types"
	orgsrepo "skema-api/features/organizations/repository"
)

type repo interface {
	Create(ctx context.Context, rec *types.EncryptedRecord) error
	FindByID(ctx context.Context, id string) (*types.EncryptedRecord, error)
	ListByOrg(ctx context.Context, orgID string) ([]*types.EncryptedRecord, error)
	Update(ctx context.Context, rec *types.EncryptedRecord) error
	Delete(ctx context.Context, id string) error
}

type Service struct {
	repo          repo
	orgsRepo      *orgsrepo.Repository
	encryptionKey string
}

func New(repo repo, orgsRepo *orgsrepo.Repository, encryptionKey string) *Service {
	return &Service{repo: repo, orgsRepo: orgsRepo, encryptionKey: encryptionKey}
}
