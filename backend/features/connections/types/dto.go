package types

// CreateConnectionRequest représente les données pour créer une connexion.
//
//	@Description	Données de création d'une connexion base de données
type CreateConnectionRequest struct {
	Name     string `json:"name"     validate:"required,min=1,max=255" example:"Production PG"`
	Driver   string `json:"driver"   validate:"required,oneof=postgres mysql" example:"postgres"`
	Host     string `json:"host"     validate:"required" example:"db.skemacms.com"`
	Port     int    `json:"port"     validate:"required,min=1,max=65535" example:"5432"`
	Database string `json:"database" validate:"required" example:"myapp"`
	User     string `json:"user"     validate:"required" example:"admin"`
	Password string `json:"password" validate:"required" example:"s3cr3t"`
	SSLMode  string `json:"ssl_mode" validate:"omitempty,oneof=disable require verify-ca verify-full" example:"disable"`
}

// UpdateConnectionRequest représente les données pour mettre à jour une connexion.
//
//	@Description	Données de mise à jour d'une connexion base de données
type UpdateConnectionRequest struct {
	Name     string `json:"name"     validate:"omitempty,min=1,max=255" example:"Production PG v2"`
	Host     string `json:"host"     validate:"omitempty" example:"db2.skemacms.com"`
	Port     int    `json:"port"     validate:"omitempty,min=1,max=65535" example:"5433"`
	Database string `json:"database" validate:"omitempty" example:"myapp_v2"`
	User     string `json:"user"     validate:"omitempty" example:"readonly"`
	Password string `json:"password" validate:"omitempty" example:"newpass"`
	SSLMode  string `json:"ssl_mode" validate:"omitempty,oneof=disable require verify-ca verify-full" example:"require"`
}

// ConnectionResponse est la représentation publique d'une connexion (sans mot de passe).
//
//	@Description	Connexion base de données (mot de passe masqué)
type ConnectionResponse struct {
	ID             string `json:"id"              example:"550e8400-e29b-41d4-a716-446655440000"`
	OrganizationID string `json:"organization_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	Name           string `json:"name"            example:"Production PG"`
	Driver         string `json:"driver"          example:"postgres"`
	Host           string `json:"host"            example:"db.skemacms.com"`
	Port           int    `json:"port"            example:"5432"`
	Database       string `json:"database"        example:"myapp"`
	User           string `json:"user"            example:"admin"`
	SSLMode        string `json:"ssl_mode"        example:"disable"`
	CreatedAt      string `json:"created_at"      example:"2025-01-01T00:00:00Z"`
}
