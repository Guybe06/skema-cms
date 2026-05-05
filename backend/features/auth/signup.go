package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"skema-api/core/mailer"
	"skema-api/features/auth/constants"
	"skema-api/features/auth/helpers"
)

/*
 * Register crée un nouveau compte utilisateur et envoie un email de vérification.
 *
 * Attend  : les informations d'inscription valides (email unique, mot de passe >= 8 chars).
 * Retourne: les tokens d'accès ou une erreur si l'email est déjà utilisé.
 */

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*TokenResponse, error) {
	existing, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New(constants.ErrEmailTaken)
	}

	hash, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &User{
		ID: uuid.NewString(), Email: req.Email, PasswordHash: hash,
		FirstName: req.FirstName, LastName: req.LastName,
		EmailVerified: false, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	if err := s.sendVerificationEmail(ctx, user); err != nil {
		return nil, err
	}

	return s.buildSession(ctx, user)
}

func (s *Service) sendVerificationEmail(ctx context.Context, user *User) error {
	raw, hashed, err := helpers.GenerateToken()
	if err != nil {
		return err
	}

	now := time.Now()
	t := &VerificationToken{
		ID: uuid.NewString(), UserID: user.ID, TokenHash: hashed,
		Type: constants.TokenTypeEmailVerification,
		ExpiresAt: now.Add(constants.VerifyTokenExpiry), CreatedAt: now,
	}
	if err := s.repo.CreateVerificationToken(ctx, t); err != nil {
		return err
	}

	return s.mailer.Send(mailer.Email{
		To:      user.Email,
		Subject: constants.SubjectVerification,
		HTML:    helpers.VerificationEmailHTML(s.frontendURL, raw),
	})
}
