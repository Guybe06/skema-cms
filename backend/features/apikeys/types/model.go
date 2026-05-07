package types

import (
	"encoding/json"
	"time"
)

type APIKey struct {
	ID                 string
	OrganizationID     string
	Name               string
	KeyHash            string
	KeyPrefix          string
	Permissions        json.RawMessage
	AllowedCollections json.RawMessage
	ExpiresAt          *time.Time
	LastUsedAt         *time.Time
	CreatedAt          time.Time
}

// Permissions représente les droits d'une clé API.
type Permissions struct {
	Read   bool `json:"read"`
	Create bool `json:"create"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}
