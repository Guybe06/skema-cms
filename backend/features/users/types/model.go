package types

import "time"

type User struct {
	ID            string
	Email         string
	FirstName     string
	LastName      string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
