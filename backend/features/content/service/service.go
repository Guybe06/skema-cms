package service

import (
	"context"

	"skema-api/core/conduit"
	colltypes "skema-api/features/collections/types"
	orgsrepo "skema-api/features/organizations/repository"
)

type collectionRepo interface {
	FindByID(ctx context.Context, id string) (*colltypes.Collection, error)
	ListFields(ctx context.Context, collectionID string) ([]*colltypes.Field, error)
}

type conduitOpener interface {
	OpenConduit(ctx context.Context, connectionID string) (conduit.Conduit, error)
}

type Service struct {
	collRepo *collectionGetter
	orgsRepo *orgsrepo.Repository
	connSvc  conduitOpener
}

// collectionGetter regroupe les méthodes de collection nécessaires au service de contenu.
type collectionGetter struct {
	repo collectionRepo
}

func New(collRepo collectionRepo, orgsRepo *orgsrepo.Repository, connSvc conduitOpener) *Service {
	return &Service{
		collRepo: &collectionGetter{repo: collRepo},
		orgsRepo: orgsRepo,
		connSvc:  connSvc,
	}
}
