package types

import "encoding/json"

// CreateAPIKeyRequest représente les données pour créer une clé API.
type CreateAPIKeyRequest struct {
	Name               string          `json:"name"                validate:"required,min=1,max=255" example:"Frontend public"`
	Permissions        Permissions     `json:"permissions"`
	AllowedCollections json.RawMessage `json:"allowed_collections" swaggertype:"array,string"`
	ExpiresAt          string          `json:"expires_at"          validate:"omitempty" example:"2026-12-31T00:00:00Z"`
}

// APIKeyResponse est la représentation publique d'une clé (sans hash, avec préfixe).
type APIKeyResponse struct {
	ID                 string          `json:"id"                   example:"550e8400-e29b-41d4-a716-446655440000"`
	OrganizationID     string          `json:"organization_id"      example:"660e8400-e29b-41d4-a716-446655440001"`
	Name               string          `json:"name"                 example:"Frontend public"`
	KeyPrefix          string          `json:"key_prefix"           example:"sk_live_1a2b"`
	Permissions        json.RawMessage `json:"permissions"          swaggertype:"object"`
	AllowedCollections json.RawMessage `json:"allowed_collections"  swaggertype:"array,string"`
	CreatedAt          string          `json:"created_at"           example:"2025-01-01T00:00:00Z"`
}

// APIKeyCreatedResponse est retourné une seule fois à la création (contient la clé brute).
type APIKeyCreatedResponse struct {
	APIKeyResponse
	RawKey string `json:"raw_key" example:"sk_live_1a2b3c4d..."`
}
