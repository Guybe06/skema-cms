package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"skema-api/core/conduit"
)

type Adapter struct {
	pool *pgxpool.Pool
}

/*
 * New ouvre une connexion PostgreSQL via pgx et vérifie la disponibilité.
 *
 * Attend  : une DSN PostgreSQL valide (postgres://user:pass@host:port/db).
 * Retourne: un Adapter prêt à l'emploi ou une erreur de connexion.
 */

func New(ctx context.Context, dsn string) (*Adapter, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = conduit.PoolMaxConns
	cfg.MinConns = conduit.PoolMinConns

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &Adapter{pool: pool}, nil
}

func (a *Adapter) Driver() string { return conduit.DriverPostgres }

func (a *Adapter) Ping(ctx context.Context) error {
	return a.pool.Ping(ctx)
}

func (a *Adapter) Close() error {
	a.pool.Close()
	return nil
}

func (a *Adapter) Exec(ctx context.Context, query string, args ...any) (conduit.Result, error) {
	tag, err := a.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &result{rows: tag.RowsAffected()}, nil
}

func (a *Adapter) QueryRow(ctx context.Context, query string, args ...any) conduit.Row {
	return a.pool.QueryRow(ctx, query, args...)
}

func (a *Adapter) Query(ctx context.Context, query string, args ...any) (conduit.Rows, error) {
	rows, err := a.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &pgRows{rows: rows}, nil
}
