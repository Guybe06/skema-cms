package service

import (
	"context"
	"errors"

	"skema-api/features/apikeys/constants"
	"skema-api/features/apikeys/types"
)

func (s *Service) List(ctx context.Context, requesterID, orgSlug string) ([]*types.APIKey, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return nil, errors.New("organisation introuvable")
	}
	if org.OwnerID != requesterID {
		return nil, errors.New(constants.ErrNotAuthorized)
	}
	return s.repo.ListByOrg(ctx, org.ID)
}

func (s *Service) Revoke(ctx context.Context, requesterID, orgSlug, id string) error {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return errors.New("organisation introuvable")
	}
	if org.OwnerID != requesterID {
		return errors.New(constants.ErrNotAuthorized)
	}

	k, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if k == nil || k.OrganizationID != org.ID {
		return errors.New(constants.ErrKeyNotFound)
	}
	return s.repo.Delete(ctx, k.ID)
}
