package types

import (
	"encoding/json"
	"time"
)

type Collection struct {
	ID             string
	ConnectionID   string
	OrganizationID string
	Name           string
	TableName      string
	DisplayName    string
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Fields         []*Field
}

type Field struct {
	ID           string
	CollectionID string
	Name         string
	ColumnName   string
	Type         string
	Required     bool
	IsUnique     bool
	DefaultValue string
	Options      json.RawMessage
	Position     int
	CreatedAt    time.Time
}
