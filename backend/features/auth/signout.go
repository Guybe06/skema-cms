package auth

import (
	"context"
	"errors"

	"skema-api/features/auth/constants"
)

/*
 * Logout invalide la session courante de l'utilisateur.
 *
 * Attend  : l'identifiant de session extrait du token JWT.
 * Retourne: une erreur si la session est introuvable.
 */

func (s *Service) Logout(ctx context.Context, sessionID string) error {
	if err := s.repo.DeleteSession(ctx, sessionID); err != nil {
		return errors.New(constants.ErrSessionNotFound)
	}
	return nil
}
