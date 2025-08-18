package main

import (
	"fmt"
	"log"
	"os"
	"social-media-go-ddd/internal/application/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Accept arg "up" or "down", default "up"
	action := "up"
	if len(os.Args) > 1 {
		action = os.Args[1]
	}

	cfg := config.LoadConfig()

	var migrationsPath string

	switch cfg.DBDriver {
	case config.DB_DRIVER_PG:
		migrationsPath = "file://internal/infrastructure/persistence/postgres/migrations"
	case config.DB_DRIVER_MYSQL:
		migrationsPath = "file://internal/infrastructure/persistence/mysql/migrations"
	default:
		log.Fatalf("Unsupported DB driver: %s", cfg.DBDriver)
	}

	// DSN compatible with golang-migrate
	var dsn string
	switch cfg.DBDriver {
	case config.DB_DRIVER_PG:
		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
		)
	case config.DB_DRIVER_MYSQL:
		dsn = fmt.Sprintf(
			"mysql://%s:%s@tcp(%s:%s)/%s",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
		)
	}

	fmt.Printf("Equivalent command: migrate -source '%s' -database '%s' '%s'\n", migrationsPath, dsn, action)

	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		log.Fatalf("Migration init failed: %v", err)
	}

	switch action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		fmt.Println("Migrations applied successfully!")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		fmt.Println("Rollback applied successfully!")
	default:
		log.Fatalf("Unsupported action: %s. Use 'up' or 'down'", action)
	}
}
