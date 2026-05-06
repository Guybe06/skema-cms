package service

import (
	"context"
	"errors"

	"skema-api/features/accounts/constants"
	"skema-api/features/accounts/helpers"
	"skema-api/features/accounts/types"
)

/*
 * Refresh valide le refresh token, effectue sa rotation et retourne de nouveaux tokens.
 *
 * Attend  : le refresh token brut reçu du client.
 * Retourne: une nouvelle paire de tokens ou une erreur si le token est invalide.
 */

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*types.TokenResponse, error) {
	tokenHash := helpers.HashToken(refreshToken)

	session, err := s.repo.FindSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New(constants.ErrSessionNotFound)
	}

	user, err := s.repo.FindUserByID(ctx, session.UserID)
	if err != nil || user == nil {
		return nil, errors.New(constants.ErrSessionNotFound)
	}

	if err := s.repo.DeleteSession(ctx, session.ID); err != nil {
		return nil, err
	}

	return s.buildSession(ctx, user)
}
