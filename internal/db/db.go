package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/redbonzai/user-management-api/internal/config"
)

var DB *sql.DB

func InitDB(cfg *config.Config) {
	var err error
	DB, err = sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("could not ping the database: %v", err)
	}

	log.Println("connected to the database successfully")
}
