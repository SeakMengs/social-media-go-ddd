package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"social-media-go-ddd/internal/application/config"
	"social-media-go-ddd/internal/application/http"
	"social-media-go-ddd/internal/application/service"
	"social-media-go-ddd/internal/domain/repository"
	"social-media-go-ddd/internal/infrastructure/persistence/postgres"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()

	if err := cfg.DBDriverValid(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	var pool *pgxpool.Pool
	var err error

	var userRepo repository.UserRepository
	var sessionRepo repository.SessionRepository

	switch cfg.DBDriver {
	case config.DB_DRIVER_PG:
		pool, err = postgres.NewPgPool(ctx, cfg.BuildDSN())
		if err != nil {
			log.Fatalf("Failed to create Postgres pool: %v", err)
		}
		defer pool.Close()

		if err = pool.Ping(ctx); err != nil {
			log.Fatalf("Failed to ping Postgres: %v", err)
		}

		userRepo = postgres.NewPgUserRepository(pool)
		sessionRepo = postgres.NewPgSessionRepository(pool)
	case config.DB_DRIVER_MYSQL:
		// userRepo = repository.NewUserRepositoryMySQL(dbPool)
	}

	userService := service.NewUserService(userRepo)
	sessionService := service.NewSessionService(sessionRepo)

	authMiddleware := http.NewAuthMiddleware(sessionService, userService)

	userHandler := http.NewUserHandler(userService, sessionService, authMiddleware)

	app := fiber.New()

	userHandler.RegisterRoutes(app)

	// Run server on another goroutine such that we can handle graceful shutdown
	go func() {
		fmt.Printf("Server running on port %s\n", cfg.AppPort)
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Printf("Fiber server stopped: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	sigCh := make(chan os.Signal, 1)
	// SIGINT means Ctrl+C, SIGTERM is a termination signal
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("\nShutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	if pool != nil {
		log.Println("Closing postgres connection pool...")
		pool.Close()
	}

	log.Println("Server gracefully stopped.")
}
