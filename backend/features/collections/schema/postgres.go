package schema

import (
	"fmt"
	"strings"

	"skema-api/features/collections/constants"
	"skema-api/features/collections/types"
)

type postgresBuilder struct{}

func (b *postgresBuilder) Driver() string { return "postgres" }

func (b *postgresBuilder) CreateTable(tableName string) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)`, pgQuote(tableName))
}

func (b *postgresBuilder) AddColumn(tableName string, f *types.Field) string {
	colType := pgColumnType(f.Type)
	col := pgQuote(f.ColumnName)
	table := pgQuote(tableName)

	parts := []string{fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, col, colType)}

	var constraints []string
	if f.Required {
		constraints = append(constraints, "NOT NULL")
	}
	if f.IsUnique {
		constraints = append(constraints, "UNIQUE")
	}
	if f.DefaultValue != "" {
		constraints = append(constraints, fmt.Sprintf("DEFAULT %s", pgDefaultValue(f.Type, f.DefaultValue)))
	}

	if len(constraints) > 0 {
		parts = append(parts, strings.Join(constraints, " "))
	}
	return strings.Join(parts, " ")
}

func (b *postgresBuilder) DropColumn(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS %s", pgQuote(tableName), pgQuote(columnName))
}

func (b *postgresBuilder) DropTable(tableName string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", pgQuote(tableName))
}

func pgQuote(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

func pgColumnType(fieldType string) string {
	switch fieldType {
	case constants.FieldTypeText, constants.FieldTypeEmail,
		constants.FieldTypeSelect, constants.FieldTypeSlug:
		return "VARCHAR(255)"
	case constants.FieldTypeURL:
		return "VARCHAR(500)"
	case constants.FieldTypePhone, constants.FieldTypeColor:
		return "VARCHAR(50)"
	case constants.FieldTypePassword, constants.FieldTypeFile,
		constants.FieldTypeImage, constants.FieldTypeTextarea,
		constants.FieldTypeRichtext:
		return "TEXT"
	case constants.FieldTypeNumber:
		return "NUMERIC(15,4)"
	case constants.FieldTypeBoolean:
		return "BOOLEAN"
	case constants.FieldTypeDate:
		return "DATE"
	case constants.FieldTypeDatetime:
		return "TIMESTAMPTZ"
	case constants.FieldTypeMultiselect:
		return "TEXT[]"
	case constants.FieldTypeJSON:
		return "JSONB"
	default:
		return "TEXT"
	}
}

func pgDefaultValue(fieldType, value string) string {
	switch fieldType {
	case constants.FieldTypeBoolean:
		return value
	case constants.FieldTypeNumber:
		return value
	default:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(value, "'", "''"))
	}
}
