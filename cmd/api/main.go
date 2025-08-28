package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"social-media-go-ddd/internal/application/config"
	"social-media-go-ddd/internal/application/http"
	"social-media-go-ddd/internal/application/service"
	"social-media-go-ddd/internal/domain/repository"
	"social-media-go-ddd/internal/infrastructure/cache"
	"social-media-go-ddd/internal/infrastructure/cache/redis"
	"social-media-go-ddd/internal/infrastructure/persistence/mysql"
	"social-media-go-ddd/internal/infrastructure/persistence/postgres"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()

	if err := cfg.DB.DriverValid(); err != nil {
		log.Fatalf("Database driver invalid error: %v", err)
	}

	var pool *pgxpool.Pool
	var mysqlDB *sql.DB

	var err error

	var userRepo repository.UserRepository
	var sessionRepo repository.SessionRepository
	var postRepo repository.PostRepository
	var favoriteRepo repository.FavoriteRepository
	var likeRepo repository.LikeRepository
	var repostRepo repository.RepostRepository
	var followRepo repository.FollowRepository

	switch cfg.DB.Driver {
	case config.DB_DRIVER_PG:
		pool, err = postgres.NewPgPool(ctx, cfg.DB.BuildDSN())
		if err != nil {
			log.Fatalf("Failed to create Postgres pool: %v", err)
		}
		defer pool.Close()
		log.Println("Postgres connection pool established")

		userRepo = postgres.NewPgUserRepository(pool)
		sessionRepo = postgres.NewPgSessionRepository(pool)
		postRepo = postgres.NewPgPostRepository(pool)
		favoriteRepo = postgres.NewPgFavoriteRepository(pool)
		likeRepo = postgres.NewPgLikeRepository(pool)
		repostRepo = postgres.NewPgRepostRepository(pool)
		followRepo = postgres.NewPgFollowRepository(pool)
	case config.DB_DRIVER_MYSQL:
		mysqlDB, err = mysql.NewMySQLDB(cfg.DB.BuildDSN())
		if err != nil {
			log.Fatalf("Failed to connect to MySQL: %v", err)
		}
		defer mysqlDB.Close()
		log.Println("MySQL connection established")

		userRepo = mysql.NewMySQLUserRepository(mysqlDB)
		sessionRepo = mysql.NewMySQLSessionRepository(mysqlDB)
		postRepo = mysql.NewMySQLPostRepository(mysqlDB)
		favoriteRepo = mysql.NewMySQLFavoriteRepository(mysqlDB)
		likeRepo = mysql.NewMySQLLikeRepository(mysqlDB)
		repostRepo = mysql.NewMySQLRepostRepository(mysqlDB)
		followRepo = mysql.NewMySQLFollowRepository(mysqlDB)
	}

	var cacheClient cache.Cache
	cacheClient, err = redis.NewRedisCache(ctx, cfg.Redis.Addr(), cfg.Redis.Password, cfg.Redis.DB, cfg.DB.Driver)
	if err != nil {
		log.Fatalf("Failed to connect Redis cache: %v", err)
	}

	userService := service.NewUserService(userRepo, cacheClient)
	sessionService := service.NewSessionService(sessionRepo, cacheClient)
	postService := service.NewPostService(postRepo, cacheClient)
	favoriteService := service.NewFavoriteService(favoriteRepo, cacheClient)
	likeService := service.NewLikeService(likeRepo, cacheClient)
	repostService := service.NewRepostService(repostRepo, cacheClient)
	followService := service.NewFollowService(followRepo, cacheClient)

	authMiddleware := http.NewAuthMiddleware(sessionService, userService)

	userHandler := http.NewUserHandler(userService, sessionService, postService, repostService, followService, favoriteService, authMiddleware)
	postHandler := http.NewPostHandler(postService, likeService, repostService, favoriteService, sessionService, authMiddleware)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())

	userHandler.RegisterRoutes(app)
	postHandler.RegisterRoutes(app)

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

	log.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	if pool != nil {
		log.Println("Closing postgres connection pool...")
		pool.Close()
	}

	if mysqlDB != nil {
		log.Println("Closing MySQL connection...")
		if err := mysqlDB.Close(); err != nil {
			log.Printf("Error closing MySQL connection: %v", err)
		}
	}

	log.Println("Closing cache connection...")
	if err := cacheClient.Close(); err != nil {
		log.Printf("Error closing cache connection: %v", err)
	}

	log.Println("Cleanup completed. Exiting application.")

	log.Println("Server gracefully stopped.")
}
