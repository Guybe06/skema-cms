package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"skema-api/core/conduit"
	"skema-api/features/collections/types"
)

/*
 * ScanRows lit toutes les lignes d'un résultat Conduit et les convertit en maps.
 *
 * Attend  : un ensemble de lignes Conduit.
 * Retourne: une liste de maps colonne→valeur, ou une erreur.
 */

func ScanRows(rows conduit.Rows) ([]map[string]any, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var result []map[string]any
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := make(map[string]any, len(cols))
		for i, col := range cols {
			row[col] = normalizeVal(vals[i])
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

/*
 * normalizeVal convertit les types pgx non sérialisables en types JSON natifs.
 *
 * Attend  : une valeur quelconque issue d'un scan pgx.
 * Retourne: une valeur sérialisable en JSON.
 */

func normalizeVal(v any) any {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case [16]byte:
		return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
			val[0:4], val[4:6], val[6:8], val[8:10], val[10:16])
	case []byte:
		return string(val)
	case int32:
		return int64(val)
	case int8:
		return int64(val)
	case float32:
		return float64(val)
	default:
		type stringer interface{ String() string }
		if s, ok := v.(stringer); ok {
			str := s.String()
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				return f
			}
			return str
		}
		return v
	}
}

func fetchByID(ctx context.Context, conn conduit.Conduit, table string, fields []*types.Field, id string) (map[string]any, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE id=$1::uuid", BuildSelectCols(fields), quoteIdent(table))
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	entries, err := ScanRows(rows)
	if err != nil || len(entries) == 0 {
		return nil, errors.New("entrée introuvable")
	}
	return entries[0], nil
}
