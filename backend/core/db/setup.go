package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"skema-api/core/config"
)

/*
 * EnsureExists vérifie que la base de données CMS existe et la crée si nécessaire.
 *
 * Attend  : la configuration de base de données issue de config.Load().
 * Retourne: une erreur si la connexion ou la création échoue.
 */

func EnsureExists(ctx context.Context, cfg config.DatabaseConfig) error {
	adminDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/postgres?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port,
	)

	conn, err := pgx.Connect(ctx, adminDSN)
	if err != nil {
		return fmt.Errorf("connexion admin impossible : %w", err)
	}
	defer conn.Close(ctx)

	var exists bool
	err = conn.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`, cfg.Name,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err = conn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.Name))
		if err != nil {
			return fmt.Errorf("création base de données échouée : %w", err)
		}
		log.Printf("Base de données '%s' créée avec succès.", cfg.Name)
	}

	return nil
}
