package schema

import (
	"fmt"
	"skema-api/features/collections/types"
)

// Builder génère les instructions DDL pour un driver donné.
type Builder interface {
	Driver() string
	CreateTable(tableName string) string
	AddColumn(tableName string, field *types.Field) string
	DropColumn(tableName, columnName string) string
	DropTable(tableName string) string
}

// New retourne le builder adapté au driver demandé.
func New(driver string) (Builder, error) {
	switch driver {
	case "postgres":
		return &postgresBuilder{}, nil
	case "mysql":
		return &mysqlBuilder{}, nil
	default:
		return nil, fmt.Errorf("driver non supporté : %s", driver)
	}
}
