package service

import (
	"context"
	"errors"

	"skema-api/features/users/constants"
	"skema-api/features/users/types"
)

func (s *Service) GetProfile(ctx context.Context, userID string) (*types.User, error) {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New(constants.ErrUserNotFound)
	}
	return u, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID, firstName, lastName string) (*types.User, error) {
	if err := s.repo.UpdateProfile(ctx, userID, firstName, lastName); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, userID)
}
