package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"skema-api/features/accounts/constants"
	"skema-api/features/accounts/helpers"
	"skema-api/features/accounts/types"
)

var timeNow = time.Now

/*
 * buildSession crée une session avec rotation de token et génère les tokens JWT.
 *
 * Attend  : un utilisateur valide et authentifié.
 * Retourne: la réponse avec les tokens d'accès et de rafraîchissement.
 */

func (s *Service) buildSession(ctx context.Context, user *types.User) (*types.TokenResponse, error) {
	rawRefresh, hashedRefresh, err := helpers.GenerateToken()
	if err != nil {
		return nil, err
	}

	sessionID := uuid.NewString()
	now := timeNow()
	session := &types.Session{
		ID:        sessionID,
		UserID:    user.ID,
		TokenHash: hashedRefresh,
		ExpiresAt: now.Add(constants.RefreshTokenExpiry),
		CreatedAt: now,
	}
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	accessToken, err := helpers.GenerateJWT(user.ID, sessionID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &types.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresIn:    int(constants.AccessTokenExpiry.Seconds()),
		User: types.UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			EmailVerified: user.EmailVerified,
		},
	}, nil
}
