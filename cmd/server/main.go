package main

import (
	"fmt"
	"log"

	"github.com/redbonzai/user-management-api/internal/config"
	"github.com/redbonzai/user-management-api/internal/db"
	"github.com/redbonzai/user-management-api/internal/infrastructure"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	fmt.Printf("Config loaded successfully %v+\n", cfg)
	db.InitDB(cfg)

	router := infrastructure.NewRouter()

	router.Logger.Fatal(router.Start(cfg.ServerAddress))
}
