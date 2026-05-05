package mysql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"skema-api/core/conduit"
)

type Adapter struct {
	db *sql.DB
}

/*
 * New ouvre une connexion MySQL via database/sql et vérifie la disponibilité.
 *
 * Attend  : une DSN MySQL valide (user:pass@tcp(host:port)/db?parseTime=true).
 * Retourne: un Adapter prêt à l'emploi ou une erreur de connexion.
 */

func New(ctx context.Context, dsn string) (*Adapter, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(int(conduit.PoolMaxConns))
	db.SetMaxIdleConns(int(conduit.PoolMinConns))

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("%s : %w", conduit.ErrPingFailed, err)
	}

	return &Adapter{db: db}, nil
}

func (a *Adapter) Driver() string { return conduit.DriverMySQL }

func (a *Adapter) Ping(ctx context.Context) error {
	return a.db.PingContext(ctx)
}

func (a *Adapter) Close() error {
	return a.db.Close()
}

func (a *Adapter) Exec(ctx context.Context, query string, args ...any) (conduit.Result, error) {
	res, err := a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	rows, _ := res.RowsAffected()
	return &result{rows: rows}, nil
}

func (a *Adapter) QueryRow(ctx context.Context, query string, args ...any) conduit.Row {
	return a.db.QueryRowContext(ctx, query, args...)
}

func (a *Adapter) Query(ctx context.Context, query string, args ...any) (conduit.Rows, error) {
	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
