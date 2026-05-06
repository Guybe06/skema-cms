package types

import "time"

type Organization struct {
	ID        string
	Name      string
	Slug      string
	OwnerID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
