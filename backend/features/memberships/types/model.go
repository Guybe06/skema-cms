package types

import "time"

type Membership struct {
	ID             string
	OrganizationID string
	UserID         string
	Email          string
	Role           string
	Status         string
	InvitedBy      string
	TokenHash      string
	ExpiresAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
