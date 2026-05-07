package service

import (
	"fmt"
	"strings"

	"skema-api/features/collections/types"
)

func quoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

// BuildSelectCols construit la liste des colonnes SELECT en castant les types problématiques en texte.
func BuildSelectCols(fields []*types.Field) string {
	parts := []string{"id::text AS id"}
	for _, f := range fields {
		col := quoteIdent(f.ColumnName)
		switch f.Type {
		case "datetime", "date":
			parts = append(parts, fmt.Sprintf("COALESCE(%s::text, '') AS %s", col, quoteIdent(f.ColumnName)))
		case "json":
			parts = append(parts, fmt.Sprintf("COALESCE(%s::text, 'null') AS %s", col, quoteIdent(f.ColumnName)))
		default:
			parts = append(parts, col)
		}
	}
	parts = append(parts, "created_at::text AS created_at", "updated_at::text AS updated_at")
	return strings.Join(parts, ", ")
}

// BuildInsert construit un INSERT dynamique et retourne la requête et ses arguments.
func BuildInsert(table string, fields []*types.Field, data map[string]any) (string, []any, error) {
	cols := []string{}
	placeholders := []string{}
	args := []any{}

	for _, f := range fields {
		val, ok := data[f.Name]
		if !ok {
			if f.Required {
				return "", nil, fmt.Errorf("champ requis manquant : %s", f.Name)
			}
			continue
		}
		cols = append(cols, quoteIdent(f.ColumnName))
		args = append(args, val)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id::text AS id",
		quoteIdent(table),
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)
	if len(cols) == 0 {
		query = fmt.Sprintf(
			"INSERT INTO %s DEFAULT VALUES RETURNING id::text AS id",
			quoteIdent(table),
		)
	}
	return query, args, nil
}

// BuildUpdate construit un UPDATE dynamique.
func BuildUpdate(table string, fields []*types.Field, data map[string]any, id string) (string, []any) {
	sets := []string{"updated_at = NOW()"}
	args := []any{}

	for _, f := range fields {
		if val, ok := data[f.Name]; ok {
			args = append(args, val)
			sets = append(sets, fmt.Sprintf("%s = $%d", quoteIdent(f.ColumnName), len(args)))
		}
	}

	args = append(args, id)
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = $%d",
		quoteIdent(table),
		strings.Join(sets, ", "),
		len(args),
	)
	return query, args
}
