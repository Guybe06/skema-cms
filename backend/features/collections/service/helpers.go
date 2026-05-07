package service

import (
	"context"
	"errors"

	"skema-api/features/collections/constants"
)

func (s *Service) checkAccess(ctx context.Context, orgID, ownerID, requesterID string) error {
	if ownerID == requesterID {
		return nil
	}
	ok, err := s.orgsRepo.IsMember(ctx, orgID, requesterID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(constants.ErrNotAuthorized)
	}
	return nil
}
