package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerAddress    string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresHost     string
	DatabaseURL      string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		ServerAddress:    os.Getenv("SERVER_ADDRESS"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
	}

	// Check if any required environment variables are missing
	var missingVars []string
	if cfg.ServerAddress == "" {
		missingVars = append(missingVars, "SERVER_ADDRESS")
	}
	if cfg.PostgresUser == "" {
		missingVars = append(missingVars, "POSTGRES_USER")
	}
	if cfg.PostgresPassword == "" {
		missingVars = append(missingVars, "POSTGRES_PASSWORD")
	}
	if cfg.PostgresDB == "" {
		missingVars = append(missingVars, "POSTGRES_DB")
	}
	if cfg.PostgresHost == "" {
		missingVars = append(missingVars, "POSTGRES_HOST")
	}

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missingVars)
	}

	// Construct the database URL from individual components if not already set
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:5432/%s?sslmode=disable",
			cfg.PostgresUser,
			cfg.PostgresPassword,
			cfg.PostgresHost,
			cfg.PostgresDB,
		)
	}

	// Print all configuration values for debugging
	fmt.Printf("Config VARS: %+v\n", cfg) // %+v prints field names and values

	return cfg, nil
}
