package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"skema-api/core/mailer"
	"skema-api/features/memberships/constants"
	"skema-api/features/memberships/helpers"
	"skema-api/features/memberships/types"
)

func (s *Service) Invite(ctx context.Context, inviterID, orgSlug, email, role string) error {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return errors.New("organisation introuvable")
	}

	if org.OwnerID != inviterID {
		canInvite, err := s.repo.IsAdminOrOwner(ctx, org.ID, inviterID)
		if err != nil {
			return err
		}
		if !canInvite {
			return errors.New(constants.ErrNotAuthorized)
		}
	}

	existing, err := s.repo.FindByOrgAndEmail(ctx, org.ID, email)
	if err != nil {
		return err
	}
	if existing != nil && existing.Status == constants.StatusActive {
		return errors.New(constants.ErrAlreadyMember)
	}

	raw, tokenHash := generateInviteToken()
	now := time.Now()
	expires := now.Add(constants.InviteTokenExpiry)

	m := &types.Membership{
		ID:             uuid.New().String(),
		OrganizationID: org.ID,
		Email:          email,
		Role:           role,
		Status:         constants.StatusPending,
		InvitedBy:      inviterID,
		TokenHash:      tokenHash,
		ExpiresAt:      &expires,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := s.repo.Create(ctx, m); err != nil {
		return err
	}

	// L'échec d'envoi d'email ne bloque pas l'invitation déjà enregistrée.
	_ = s.mailer.Send(mailer.Email{
		To:      email,
		Subject: "Invitation à rejoindre " + org.Name,
		HTML:    helpers.InvitationEmailHTML(org.Name, inviterID, s.frontendURL, raw),
	})

	return nil
}

func (s *Service) AcceptInvite(ctx context.Context, userID, rawToken string) error {
	sum := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(sum[:])

	m, err := s.repo.FindByToken(ctx, tokenHash)
	if err != nil {
		return err
	}
	if m == nil || m.Status != constants.StatusPending {
		return errors.New(constants.ErrInviteTokenInvalid)
	}
	if m.ExpiresAt != nil && time.Now().After(*m.ExpiresAt) {
		return errors.New(constants.ErrInviteTokenInvalid)
	}

	return s.repo.AcceptInvite(ctx, m.ID, userID)
}

func generateInviteToken() (raw, hashed string) {
	b := make([]byte, 32)
	rand.Read(b)
	raw = hex.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	hashed = hex.EncodeToString(sum[:])
	return
}
