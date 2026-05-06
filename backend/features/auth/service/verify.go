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
 * VerifyEmail active le compte utilisateur via le token de vérification.
 *
 * Attend  : le token brut reçu par l'utilisateur dans son email.
 * Retourne: une erreur si le token est invalide ou expiré.
 */

func (s *Service) VerifyEmail(ctx context.Context, rawToken string) error {
	t, err := s.repo.FindVerificationToken(ctx, helpers.HashToken(rawToken), constants.TokenTypeEmailVerification)
	if err != nil {
		return err
	}
	if t == nil {
		return errors.New(constants.ErrTokenInvalid)
	}

	if err := s.repo.VerifyUserEmail(ctx, t.UserID); err != nil {
		return err
	}

	return s.repo.DeleteVerificationToken(ctx, t.ID)
}

/*
 * ResendVerification renvoie un email de vérification si aucun n'a été envoyé récemment.
 *
 * Attend  : l'identifiant de l'utilisateur authentifié.
 * Retourne: une erreur si un email a déjà été envoyé dans les 2 dernières minutes.
 */

func (s *Service) ResendVerification(ctx context.Context, userID string) error {
	cacheKey := fmt.Sprintf(constants.CacheKeyResend, userID)
	var cached string
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil && cached != "" {
		return errors.New(constants.ErrTokenAlreadySent)
	}

	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New(constants.ErrTokenInvalid)
	}

	_ = s.repo.DeleteTokensByUser(ctx, userID, constants.TokenTypeEmailVerification)

	raw, hashed, err := helpers.GenerateToken()
	if err != nil {
		return err
	}

	now := time.Now()
	t := &types.VerificationToken{
		ID: uuid.NewString(), UserID: userID, TokenHash: hashed,
		Type:      constants.TokenTypeEmailVerification,
		ExpiresAt: now.Add(constants.VerifyTokenExpiry), CreatedAt: now,
	}
	if err := s.repo.CreateVerificationToken(ctx, t); err != nil {
		return err
	}

	_ = s.cache.Set(ctx, cacheKey, "1", constants.CacheTTLResend)

	return s.mailer.Send(mailer.Email{
		To:      user.Email,
		Subject: constants.SubjectVerification,
		HTML:    helpers.VerificationEmailHTML(s.frontendURL, raw),
	})
}
