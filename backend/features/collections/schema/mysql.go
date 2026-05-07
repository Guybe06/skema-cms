package schema

import (
	"fmt"
	"strings"

	"skema-api/features/collections/constants"
	"skema-api/features/collections/types"
)

type mysqlBuilder struct{}

func (b *mysqlBuilder) Driver() string { return "mysql" }

func (b *mysqlBuilder) CreateTable(tableName string) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	id CHAR(36) PRIMARY KEY,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)`, myQuote(tableName))
}

func (b *mysqlBuilder) AddColumn(tableName string, f *types.Field) string {
	colType := myColumnType(f.Type)
	col := myQuote(f.ColumnName)
	table := myQuote(tableName)

	var constraints []string
	if f.Required {
		constraints = append(constraints, "NOT NULL")
	} else {
		constraints = append(constraints, "NULL")
	}
	if f.DefaultValue != "" {
		constraints = append(constraints, fmt.Sprintf("DEFAULT %s", myDefaultValue(f.Type, f.DefaultValue)))
	}
	if f.IsUnique {
		constraints = append(constraints, "UNIQUE")
	}

	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s %s",
		table, col, colType, strings.Join(constraints, " "))
}

func (b *mysqlBuilder) DropColumn(tableName, columnName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", myQuote(tableName), myQuote(columnName))
}

func (b *mysqlBuilder) DropTable(tableName string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", myQuote(tableName))
}

func myQuote(s string) string {
	return "`" + strings.ReplaceAll(s, "`", "``") + "`"
}

func myColumnType(fieldType string) string {
	switch fieldType {
	case constants.FieldTypeText, constants.FieldTypeEmail,
		constants.FieldTypeSelect, constants.FieldTypeSlug:
		return "VARCHAR(255)"
	case constants.FieldTypeURL:
		return "VARCHAR(500)"
	case constants.FieldTypePhone, constants.FieldTypeColor:
		return "VARCHAR(50)"
	case constants.FieldTypePassword, constants.FieldTypeFile,
		constants.FieldTypeImage, constants.FieldTypeTextarea:
		return "TEXT"
	case constants.FieldTypeRichtext:
		return "LONGTEXT"
	case constants.FieldTypeNumber:
		return "DECIMAL(15,4)"
	case constants.FieldTypeBoolean:
		return "TINYINT(1)"
	case constants.FieldTypeDate:
		return "DATE"
	case constants.FieldTypeDatetime:
		return "DATETIME"
	case constants.FieldTypeMultiselect, constants.FieldTypeJSON:
		return "JSON"
	default:
		return "TEXT"
	}
}

func myDefaultValue(fieldType, value string) string {
	switch fieldType {
	case constants.FieldTypeBoolean, constants.FieldTypeNumber:
		return value
	default:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(value, "'", "''"))
	}
}
