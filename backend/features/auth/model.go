package auth

import "time"

type User struct {
	ID            string
	Email         string
	PasswordHash  string
	FirstName     string
	LastName      string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Session struct {
	ID           string
	UserID       string
	TokenHash    string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type VerificationToken struct {
	ID        string
	UserID    string
	TokenHash string
	Type      string
	ExpiresAt time.Time
	CreatedAt time.Time
}
