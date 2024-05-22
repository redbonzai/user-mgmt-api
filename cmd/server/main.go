package main

import (
	"github.com/redbonzai/user-management-api/internal/config"
	"github.com/redbonzai/user-management-api/internal/db"
	"github.com/redbonzai/user-management-api/internal/infrastructure"
	"github.com/redbonzai/user-management-api/pkg/logger"
	"go.uber.org/zap"
)

// @title User Management API
// @version 1.0
// @description This is a sample server User Management server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Initialize logger
	logger.InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("could not load config:", zap.Error(err))
	}

	logger.Info("-- Config loaded successfully --")
	db.InitDB(cfg)

	router := infrastructure.NewRouter()

	// Serve static files for Swagger
	router.Static("/swagger", "docs")

	if err := router.Start(cfg.ServerAddress); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
