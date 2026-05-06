package types

import "time"

// Connection représente une connexion en mémoire (champs déchiffrés).
type Connection struct {
	ID             string
	OrganizationID string
	Name           string
	Driver         string
	Host           string
	Port           int
	Database       string
	User           string
	SSLMode        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// EncryptedRecord est le format de stockage en base (champs sensibles chiffrés).
type EncryptedRecord struct {
	ID                string
	OrganizationID    string
	Name              string
	Driver            string
	HostEncrypted     string
	PortEncrypted     string
	DatabaseEncrypted string
	UserEncrypted     string
	PasswordEncrypted string
	SSLMode           string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
