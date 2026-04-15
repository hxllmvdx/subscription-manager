package main

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/migrate"
	"backend/internal/repository"
	"backend/internal/routes"
	"backend/internal/service"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// @title Subscription Manager API
	// @version 1.0
	// @description API for managing user online subscriptions
	// @host localhost:8080
	// @BasePath /
	_ "backend/docs"
)

// @title Subscription Manager API
// @version 1.0
// @description API for managing user online subscriptions
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.LoadConfig()

	db, err := repository.NewDB(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to get sql.DB", "error", err)
		return
	}
	defer sqlDB.Close()

	if err := migrate.RunMigrations(cfg, sqlDB); err != nil {
		slog.Error("failed to run migrations", "error", err)
		return
	}

	subscriptionRepo := repository.NewSubscriptionRepo(db)

	subscriptionService := service.NewSubscriptionService(subscriptionRepo)

	h := handler.NewSubscriptionHandler(subscriptionService)

	router := gin.Default()

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := router.Group("/api")
	apiGroup.Use(middleware.RequestLogger())
	routes.SetupApiRoutes(apiGroup, h)

	srv := &http.Server{
		Addr:    cfg.BackendPort,
		Handler: router,
	}

	slog.Info("server starting", "port", cfg.BackendPort)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", "error", err)
	}
}
