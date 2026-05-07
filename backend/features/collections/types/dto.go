package types

import "encoding/json"

// CreateCollectionRequest représente les données pour créer une collection.
type CreateCollectionRequest struct {
	ConnectionID string `json:"connection_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name         string `json:"name"          validate:"required,min=1,max=255" example:"Products"`
	TableName    string `json:"table_name"    validate:"required,min=1,max=63,alphanum_underscore" example:"products"`
	DisplayName  string `json:"display_name"  validate:"omitempty,max=255" example:"Catalogue produits"`
	Description  string `json:"description"   validate:"omitempty,max=1000" example:"Liste des produits du catalogue"`
}

// UpdateCollectionRequest représente les données pour modifier une collection.
type UpdateCollectionRequest struct {
	Name        string `json:"name"         validate:"omitempty,min=1,max=255" example:"Products v2"`
	DisplayName string `json:"display_name" validate:"omitempty,max=255" example:"Catalogue"`
	Description string `json:"description"  validate:"omitempty,max=1000" example:"Description mise à jour"`
}

// AddFieldRequest représente les données pour ajouter un champ à une collection.
type AddFieldRequest struct {
	Name         string          `json:"name"          validate:"required,min=1,max=255" example:"price"`
	ColumnName   string          `json:"column_name"   validate:"required,min=1,max=63" example:"price"`
	Type         string          `json:"type"          validate:"required" example:"number"`
	Required     bool            `json:"required"      example:"false"`
	IsUnique     bool            `json:"is_unique"     example:"false"`
	DefaultValue string          `json:"default_value" validate:"omitempty" example:"0"`
	Options      json.RawMessage `json:"options"       swaggertype:"object"`
	Position     int             `json:"position"      example:"1"`
}

// CollectionResponse est la représentation publique d'une collection.
type CollectionResponse struct {
	ID             string          `json:"id"              example:"550e8400-e29b-41d4-a716-446655440000"`
	ConnectionID   string          `json:"connection_id"   example:"660e8400-e29b-41d4-a716-446655440001"`
	OrganizationID string          `json:"organization_id" example:"770e8400-e29b-41d4-a716-446655440002"`
	Name           string          `json:"name"            example:"Products"`
	TableName      string          `json:"table_name"      example:"products"`
	DisplayName    string          `json:"display_name"    example:"Catalogue produits"`
	Description    string          `json:"description"     example:"Liste des produits"`
	CreatedAt      string          `json:"created_at"      example:"2025-01-01T00:00:00Z"`
	Fields         []*FieldResponse `json:"fields,omitempty"`
}

// FieldResponse est la représentation publique d'un champ.
type FieldResponse struct {
	ID           string          `json:"id"            example:"880e8400-e29b-41d4-a716-446655440003"`
	Name         string          `json:"name"          example:"price"`
	ColumnName   string          `json:"column_name"   example:"price"`
	Type         string          `json:"type"          example:"number"`
	Required     bool            `json:"required"      example:"false"`
	IsUnique     bool            `json:"is_unique"     example:"false"`
	DefaultValue string          `json:"default_value" example:"0"`
	Options      json.RawMessage `json:"options"       swaggertype:"object"`
	Position     int             `json:"position"      example:"1"`
}
