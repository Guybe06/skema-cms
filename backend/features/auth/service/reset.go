package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"skema-api/core/mailer"
	"skema-api/features/auth/constants"
	"skema-api/features/auth/helpers"
	"skema-api/features/auth/types"
)

/*
 * RequestReset envoie un email de réinitialisation si le compte existe.
 *
 * Attend  : l'adresse email de l'utilisateur.
 * Retourne: toujours nil (pas de divulgation d'existence de compte).
 */

func (s *Service) RequestReset(ctx context.Context, email string) error {
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil || user == nil {
		return nil
	}

	cacheKey := fmt.Sprintf(constants.CacheKeyReset, user.ID)
	var cached string
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil && cached != "" {
		return nil
	}

	_ = s.repo.DeleteTokensByUser(ctx, user.ID, constants.TokenTypePasswordReset)

	raw, hashed, err := helpers.GenerateToken()
	if err != nil {
		return err
	}

	now := time.Now()
	t := &types.VerificationToken{
		ID: uuid.NewString(), UserID: user.ID, TokenHash: hashed,
		Type:      constants.TokenTypePasswordReset,
		ExpiresAt: now.Add(constants.ResetTokenExpiry), CreatedAt: now,
	}
	if err := s.repo.CreateVerificationToken(ctx, t); err != nil {
		return err
	}

	_ = s.cache.Set(ctx, cacheKey, "1", constants.CacheTTLResend)

	return s.mailer.Send(mailer.Email{
		To:      user.Email,
		Subject: constants.SubjectReset,
		HTML:    helpers.ResetEmailHTML(s.frontendURL, raw),
	})
}

/*
 * ConfirmReset applique le nouveau mot de passe après vérification du token.
 *
 * Attend  : le token brut de réinitialisation et le nouveau mot de passe.
 * Retourne: une erreur si le token est invalide ou si le hachage échoue.
 */

func (s *Service) ConfirmReset(ctx context.Context, rawToken, newPassword string) error {
	t, err := s.repo.FindVerificationToken(ctx, helpers.HashToken(rawToken), constants.TokenTypePasswordReset)
	if err != nil {
		return err
	}
	if t == nil {
		return errors.New(constants.ErrTokenInvalid)
	}

	hash, err := helpers.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateUserPassword(ctx, t.UserID, hash); err != nil {
		return err
	}

	_ = s.repo.DeleteVerificationToken(ctx, t.ID)
	_ = s.repo.DeleteUserSessions(ctx, t.UserID)
	return nil
}
