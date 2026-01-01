package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/ximofam/user-service/internal/config"
	mydb "github.com/ximofam/user-service/internal/db"
	"github.com/ximofam/user-service/internal/transport/http/middleware"
	"github.com/ximofam/user-service/internal/transport/http/routes"
	"github.com/ximofam/user-service/pkg/cache"
	"github.com/ximofam/user-service/pkg/mail"
	"github.com/ximofam/user-service/pkg/token"
)

type App struct {
	cfg         *config.Config
	db          *sql.DB
	redisClient *redis.Client
	router      *gin.Engine
}

func NewApp(cfg *config.Config, db *sql.DB, redisClient *redis.Client) (*App, error) {
	router := gin.Default()
	router.Use(
		middleware.APIKeyAuth(cfg.Server.APIKey),
	)

	dbProvider := mydb.NewDBProvider(db)
	cacheService := cache.NewCacheService(redisClient)
	tokenService := token.NewTokenService(
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
		cfg.JWT.SecretKey,
		cacheService,
	)
	mailService := mail.NewMailtrapSender(
		cfg.Mailtrap.Host,
		cfg.Mailtrap.Port,
		cfg.Mailtrap.Username,
		cfg.Mailtrap.Password,
	)

	routes.AuthRoutes(router, dbProvider, tokenService, cacheService, mailService)
	routes.UserRoutes(router, dbProvider, tokenService)

	return &App{
		cfg:         cfg,
		db:          db,
		redisClient: redisClient,
		router:      router,
	}, nil
}

func (a *App) Run() error {
	srv := http.Server{
		Addr:    a.cfg.Server.Addr,
		Handler: a.router,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Server is running at: %s", a.cfg.Server.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case <-shutdown:
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}

		log.Println("Server exited gracefully")
		return nil
	}
}
