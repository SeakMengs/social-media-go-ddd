package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	// fmt.Printf("dsn: %s\n", dsn)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Postgres DSN: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping Postgres: %v", err)
	}

	return pool, nil
}
