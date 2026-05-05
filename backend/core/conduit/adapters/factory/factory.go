package factory

import (
	"context"
	"fmt"

	"skema-api/core/conduit"
	"skema-api/core/conduit/adapters/mysql"
	"skema-api/core/conduit/adapters/postgres"
)

/*
 * New instancie l'adaptateur Conduit selon le driver demandé.
 *
 * Attend  : un driver ("postgres" ou "mysql") et une DSN de connexion.
 * Retourne: un Conduit prêt à l'emploi ou une erreur si le driver est inconnu.
 */

func New(ctx context.Context, driver, dsn string) (conduit.Conduit, error) {
	switch driver {
	case conduit.DriverPostgres:
		return postgres.New(ctx, dsn)
	case conduit.DriverMySQL:
		return mysql.New(ctx, dsn)
	default:
		return nil, fmt.Errorf(conduit.ErrUnsupportedDriver, driver)
	}
}
