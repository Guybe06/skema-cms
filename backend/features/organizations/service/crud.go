package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"skema-api/features/organizations/constants"
	"skema-api/features/organizations/types"
)

func (s *Service) Create(ctx context.Context, ownerID, name string) (*types.Organization, error) {
	slug, err := s.buildUniqueSlug(ctx, name)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	org := &types.Organization{
		ID:        uuid.New().String(),
		Name:      name,
		Slug:      slug,
		OwnerID:   ownerID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, org); err != nil {
		return nil, err
	}
	return org, nil
}

func (s *Service) ListByOwner(ctx context.Context, ownerID string) ([]*types.Organization, error) {
	return s.repo.FindByOwner(ctx, ownerID)
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*types.Organization, error) {
	org, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, errors.New(constants.ErrOrgNotFound)
	}
	return org, nil
}

func (s *Service) Update(ctx context.Context, userID, slug, name string) (*types.Organization, error) {
	org, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, errors.New(constants.ErrOrgNotFound)
	}
	if org.OwnerID != userID {
		return nil, errors.New(constants.ErrNotOwner)
	}
	newSlug, err := s.buildUniqueSlug(ctx, name)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, org.ID, name, newSlug); err != nil {
		return nil, err
	}
	org.Name = name
	org.Slug = newSlug
	return org, nil
}

func (s *Service) Delete(ctx context.Context, userID, slug string) error {
	org, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return err
	}
	if org == nil {
		return errors.New(constants.ErrOrgNotFound)
	}
	if org.OwnerID != userID {
		return errors.New(constants.ErrNotOwner)
	}
	return s.repo.Delete(ctx, org.ID)
}

func (s *Service) TransferOwnership(ctx context.Context, userID, slug, newOwnerID string) error {
	org, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return err
	}
	if org == nil {
		return errors.New(constants.ErrOrgNotFound)
	}
	if org.OwnerID != userID {
		return errors.New(constants.ErrNotOwner)
	}
	isMember, err := s.repo.IsMember(ctx, org.ID, newOwnerID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New(constants.ErrNewOwnerNotMember)
	}
	return s.repo.TransferOwnership(ctx, org.ID, newOwnerID)
}
