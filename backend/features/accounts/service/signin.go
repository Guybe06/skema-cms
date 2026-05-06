package service

import (
	"context"
	"errors"

	"skema-api/features/accounts/constants"
	"skema-api/features/accounts/helpers"
	"skema-api/features/accounts/types"
)

/*
 * Login authentifie un utilisateur et crée une nouvelle session.
 *
 * Attend  : les identifiants de connexion (email + mot de passe).
 * Retourne: les tokens d'accès ou une erreur d'authentification.
 */

func (s *Service) Login(ctx context.Context, req types.LoginRequest) (*types.TokenResponse, error) {
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		helpers.HashDummy()
		return nil, errors.New(constants.ErrInvalidCredentials)
	}

	if !helpers.CheckPassword(user.PasswordHash, req.Password) {
		return nil, errors.New(constants.ErrInvalidCredentials)
	}

	return s.buildSession(ctx, user)
}
