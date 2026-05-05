package conduit

import "context"

/*
 * Conduit est l'interface d'abstraction de base de données client.
 *
 * Attend  : un contexte et les paramètres propres à chaque opération.
 * Retourne: les données ou une erreur selon l'opération exécutée.
 */

type Conduit interface {
	Ping(ctx context.Context) error
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
	Exec(ctx context.Context, query string, args ...any) (Result, error)
	Close() error
	Driver() string
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Columns() ([]string, error)
	Close() error
	Err() error
}

type Row interface {
	Scan(dest ...any) error
}

type Result interface {
	RowsAffected() int64
}
