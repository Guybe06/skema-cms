package service

import (
	"context"
	"errors"

	"skema-api/features/collections/constants"
	"skema-api/features/collections/types"
)

func (s *Service) List(ctx context.Context, requesterID, orgSlug string) ([]*types.Collection, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return nil, errors.New(constants.ErrOrgNotFound)
	}
	if err := s.checkAccess(ctx, org.ID, org.OwnerID, requesterID); err != nil {
		return nil, err
	}
	return s.repo.ListByOrg(ctx, org.ID)
}

func (s *Service) Get(ctx context.Context, requesterID, orgSlug, id string) (*types.Collection, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return nil, errors.New(constants.ErrOrgNotFound)
	}
	if err := s.checkAccess(ctx, org.ID, org.OwnerID, requesterID); err != nil {
		return nil, err
	}

	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil || c.OrganizationID != org.ID {
		return nil, errors.New(constants.ErrCollectionNotFound)
	}

	c.Fields, err = s.repo.ListFields(ctx, c.ID)
	return c, err
}
