package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"skema-api/features/apikeys/constants"
	"skema-api/features/apikeys/types"
	orgsrepo "skema-api/features/organizations/repository"
)

type repo interface {
	Create(ctx context.Context, k *types.APIKey) error
	FindByHash(ctx context.Context, hash string) (*types.APIKey, error)
	FindByID(ctx context.Context, id string) (*types.APIKey, error)
	ListByOrg(ctx context.Context, orgID string) ([]*types.APIKey, error)
	Delete(ctx context.Context, id string) error
	TouchLastUsed(ctx context.Context, id string)
}

type Service struct {
	repo     repo
	orgsRepo *orgsrepo.Repository
}

func New(repo repo, orgsRepo *orgsrepo.Repository) *Service {
	return &Service{repo: repo, orgsRepo: orgsRepo}
}

func (s *Service) Generate(ctx context.Context, requesterID, orgSlug string, req types.CreateAPIKeyRequest) (string, *types.APIKey, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return "", nil, errors.New("organisation introuvable")
	}
	if org.OwnerID != requesterID {
		return "", nil, errors.New(constants.ErrNotAuthorized)
	}

	raw, hash, prefix := generateKey()

	permsJSON, _ := json.Marshal(req.Permissions)

	k := &types.APIKey{
		ID:                 uuid.New().String(),
		OrganizationID:     org.ID,
		Name:               req.Name,
		KeyHash:            hash,
		KeyPrefix:          prefix,
		Permissions:        permsJSON,
		AllowedCollections: req.AllowedCollections,
		CreatedAt:          time.Now(),
	}

	if req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err == nil {
			k.ExpiresAt = &t
		}
	}

	if err := s.repo.Create(ctx, k); err != nil {
		return "", nil, err
	}
	return raw, k, nil
}

func (s *Service) List(ctx context.Context, requesterID, orgSlug string) ([]*types.APIKey, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return nil, errors.New("organisation introuvable")
	}
	if org.OwnerID != requesterID {
		return nil, errors.New(constants.ErrNotAuthorized)
	}
	return s.repo.ListByOrg(ctx, org.ID)
}

func (s *Service) Revoke(ctx context.Context, requesterID, orgSlug, id string) error {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return errors.New("organisation introuvable")
	}
	if org.OwnerID != requesterID {
		return errors.New(constants.ErrNotAuthorized)
	}

	k, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if k == nil || k.OrganizationID != org.ID {
		return errors.New(constants.ErrKeyNotFound)
	}
	return s.repo.Delete(ctx, k.ID)
}

func (s *Service) FindByHash(ctx context.Context, hash string) (*types.APIKey, error) {
	return s.repo.FindByHash(ctx, hash)
}

func (s *Service) TouchLastUsed(ctx context.Context, id string) {
	s.repo.TouchLastUsed(ctx, id)
}

// generateKey produit une clé brute, son hash SHA-256 et son préfixe d'affichage.
func generateKey() (raw, hash, prefix string) {
	b := make([]byte, 32)
	rand.Read(b)
	raw = "sk_live_" + hex.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	hash = hex.EncodeToString(sum[:])
	if len(raw) >= 12 {
		prefix = raw[:12]
	} else {
		prefix = raw
	}
	return
}
