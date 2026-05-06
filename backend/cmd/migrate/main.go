package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	steps := flag.Int("steps", 0, "number of migrations to run (0 = all)")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("usage: migrate [up|down] [-steps N]")
	}
	direction := args[0]

	ensureDB()

	dsn := buildDSN()
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalf("migrate init: %v", err)
	}
	defer m.Close()

	switch direction {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-(*steps))
		} else {
			err = m.Down()
		}
	case "force":
		if *steps == 0 {
			log.Fatal("force requires -steps <version>")
		}
		err = m.Force(*steps)
	default:
		log.Fatalf("unknown direction: %s (use up|down|force)", direction)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}

	v, dirty, _ := m.Version()
	fmt.Printf("migration %s done — version: %d, dirty: %v\n", direction, v, dirty)
}

func ensureDB() {
	host := getenv("CMS_DB_HOST", "localhost")
	port := getenv("CMS_DB_PORT", "5432")
	user := getenv("CMS_DB_USER", "postgres")
	pass := getenv("CMS_DB_PASSWORD", "")
	name := getenv("CMS_DB_NAME", "skemacms")

	adminDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", user, pass, host, port)
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, adminDSN)
	if err != nil {
		log.Fatalf("connect to postgres: %v", err)
	}
	defer conn.Close(ctx)

	var exists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", name).Scan(&exists)
	if err != nil {
		log.Fatalf("check db exists: %v", err)
	}
	if !exists {
		_, err = conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", name))
		if err != nil {
			log.Fatalf("create database: %v", err)
		}
		fmt.Printf("database %q created\n", name)
	}
}

func buildDSN() string {
	host := getenv("CMS_DB_HOST", "localhost")
	port := getenv("CMS_DB_PORT", "5432")
	name := getenv("CMS_DB_NAME", "skemacms")
	user := getenv("CMS_DB_USER", "postgres")
	pass := getenv("CMS_DB_PASSWORD", "")
	return fmt.Sprintf("pgx5://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
