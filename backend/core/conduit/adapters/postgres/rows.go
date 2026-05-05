package postgres

import "github.com/jackc/pgx/v5"

type pgRows struct {
	rows pgx.Rows
}

func (r *pgRows) Next() bool          { return r.rows.Next() }
func (r *pgRows) Scan(dest ...any) error { return r.rows.Scan(dest...) }
func (r *pgRows) Close() error        { r.rows.Close(); return nil }
func (r *pgRows) Err() error          { return r.rows.Err() }

func (r *pgRows) Columns() ([]string, error) {
	fields := r.rows.FieldDescriptions()
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = string(f.Name)
	}
	return names, nil
}

type result struct {
	rows int64
}

func (r *result) RowsAffected() int64 { return r.rows }
