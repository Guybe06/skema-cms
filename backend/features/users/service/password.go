package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"skema-api/features/users/constants"
)

func (s *Service) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	hash, err := s.repo.FindPasswordHash(ctx, userID)
	if err != nil {
		return err
	}
	if hash == "" {
		return errors.New(constants.ErrUserNotFound)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(currentPassword)); err != nil {
		return errors.New(constants.ErrInvalidPassword)
	}
	if currentPassword == newPassword {
		return errors.New(constants.ErrSamePassword)
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, userID, string(newHash))
}
