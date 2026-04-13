package main

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/migrate"
	"backend/internal/repository"
	"backend/internal/routes"
	"log/slog"

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

	if err := migrate.RunMigrations(cfg, sqlDB); err != nil {
		slog.Error("failed to run migrations", "error", err)
		return
	}

	subscriptionRepo := repository.NewSubscriptionRepo(db)

	h := handler.NewSubscriptionHandler(subscriptionRepo)

	router := gin.Default()

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := router.Group("/api")
	routes.SetupApiRoutes(apiGroup, h)

	slog.Info("server starting", "port", cfg.BackendPort)
	if err := router.Run(cfg.BackendPort); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
