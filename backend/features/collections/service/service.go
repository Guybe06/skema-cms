package service

import (
	"context"

	"skema-api/core/conduit"
	"skema-api/features/collections/types"
	orgsrepo "skema-api/features/organizations/repository"
)

type repo interface {
	Create(ctx context.Context, c *types.Collection) error
	FindByID(ctx context.Context, id string) (*types.Collection, error)
	FindByOrgAndTable(ctx context.Context, orgID, tableName string) (*types.Collection, error)
	ListByOrg(ctx context.Context, orgID string) ([]*types.Collection, error)
	Update(ctx context.Context, id, name, displayName, description string) error
	Delete(ctx context.Context, id string) error
	TableExists(ctx context.Context, connectionID, tableName string) (bool, error)
	AddField(ctx context.Context, f *types.Field) error
	ListFields(ctx context.Context, collectionID string) ([]*types.Field, error)
	FindField(ctx context.Context, id string) (*types.Field, error)
	ColumnExists(ctx context.Context, collectionID, columnName string) (bool, error)
	DeleteField(ctx context.Context, id string) error
}

// conduitOpener ouvre une connexion Conduit vers une base client.
type conduitOpener interface {
	OpenConduit(ctx context.Context, connectionID string) (conduit.Conduit, error)
}

type Service struct {
	repo     repo
	orgsRepo *orgsrepo.Repository
	connSvc  conduitOpener
}

func New(repo repo, orgsRepo *orgsrepo.Repository, connSvc conduitOpener) *Service {
	return &Service{repo: repo, orgsRepo: orgsRepo, connSvc: connSvc}
}
