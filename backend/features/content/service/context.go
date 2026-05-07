package service

import (
	"context"
	"errors"

	"skema-api/core/conduit"
	colltypes "skema-api/features/collections/types"
	"skema-api/features/content/constants"
)

/*
 * getContext résout l'organisation, la collection, les champs et ouvre la connexion Conduit.
 *
 * Attend  : l'ID du demandeur, le slug de l'organisation et l'ID de la collection.
 * Retourne: la collection, la connexion ouverte, les champs ou une erreur.
 */

func (s *Service) getContext(ctx context.Context, requesterID, orgSlug, collectionID string) (*colltypes.Collection, conduit.Conduit, []*colltypes.Field, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return nil, nil, nil, errors.New(constants.ErrOrgNotFound)
	}
	if org.OwnerID != requesterID {
		ok, err := s.orgsRepo.IsMember(ctx, org.ID, requesterID)
		if err != nil {
			return nil, nil, nil, err
		}
		if !ok {
			return nil, nil, nil, errors.New(constants.ErrNotAuthorized)
		}
	}

	c, err := s.collRepo.repo.FindByID(ctx, collectionID)
	if err != nil || c == nil || c.OrganizationID != org.ID {
		return nil, nil, nil, errors.New(constants.ErrCollectionNotFound)
	}

	fields, err := s.collRepo.repo.ListFields(ctx, c.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	conn, err := s.connSvc.OpenConduit(ctx, c.ConnectionID)
	if err != nil {
		return nil, nil, nil, errors.New(constants.ErrConnectionFailed)
	}
	return c, conn, fields, nil
}
