package main

import (
	"fmt"

	"github.com/redbonzai/user-management-api/internal/config"
	"github.com/redbonzai/user-management-api/internal/db"
	"github.com/redbonzai/user-management-api/internal/infrastructure"
	"github.com/redbonzai/user-management-api/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger.InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("could not load config:", zap.Error(err))
	}

	fmt.Printf("Config loaded successfully %v+\n", cfg)
	db.InitDB(cfg)

	router := infrastructure.NewRouter()

	router.Logger.Fatal(router.Start(cfg.ServerAddress))
}
