package service

import (
	"context"
	"errors"

	"skema-api/features/memberships/constants"
	"skema-api/features/memberships/types"
)

func (s *Service) ListMembers(ctx context.Context, orgSlug, requesterID string) ([]*types.Membership, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return nil, errors.New("organisation introuvable")
	}

	if org.OwnerID != requesterID {
		isMember, err := s.repo.IsMember(ctx, org.ID, requesterID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, errors.New(constants.ErrNotAuthorized)
		}
	}

	return s.repo.ListByOrg(ctx, org.ID)
}

func (s *Service) UpdateRole(ctx context.Context, requesterID, orgSlug, targetUserID, role string) error {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return errors.New("organisation introuvable")
	}
	if org.OwnerID != requesterID {
		return errors.New(constants.ErrNotAuthorized)
	}
	if targetUserID == org.OwnerID {
		return errors.New(constants.ErrCannotChangeOwner)
	}

	m, err := s.repo.FindByOrgAndUser(ctx, org.ID, targetUserID)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New(constants.ErrMemberNotFound)
	}

	return s.repo.UpdateRole(ctx, m.ID, role)
}

func (s *Service) RemoveMember(ctx context.Context, requesterID, orgSlug, targetUserID string) error {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return errors.New("organisation introuvable")
	}

	isSelf := requesterID == targetUserID
	isOwner := org.OwnerID == requesterID

	if !isSelf && !isOwner {
		isAdmin, err := s.repo.IsAdminOrOwner(ctx, org.ID, requesterID)
		if err != nil {
			return err
		}
		if !isAdmin {
			return errors.New(constants.ErrNotAuthorized)
		}
	}

	if targetUserID == org.OwnerID {
		return errors.New(constants.ErrCannotRemoveOwner)
	}

	m, err := s.repo.FindByOrgAndUser(ctx, org.ID, targetUserID)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New(constants.ErrMemberNotFound)
	}

	return s.repo.Delete(ctx, m.ID)
}
