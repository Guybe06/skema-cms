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
	"skema-api/features/apikeys/types"
)

/*
 * Generate crée une nouvelle clé API pour l'organisation donnée.
 *
 * Attend  : l'ID du demandeur, le slug de l'organisation et les paramètres de la clé.
 * Retourne: la clé brute (affichée une seule fois), la clé persistée, ou une erreur.
 */

func (s *Service) Generate(ctx context.Context, requesterID, orgSlug string, req types.CreateAPIKeyRequest) (string, *types.APIKey, error) {
	org, err := s.orgsRepo.FindBySlug(ctx, orgSlug)
	if err != nil || org == nil {
		return "", nil, errors.New("organisation introuvable")
	}
	if org.OwnerID != requesterID {
		return "", nil, errors.New("accès non autorisé")
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

/*
 * generateKey produit une clé brute sk_live_, son hash SHA-256 et son préfixe d'affichage.
 *
 * Retourne: la clé brute, le hash hexadécimal, les 12 premiers caractères de la clé.
 */

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
