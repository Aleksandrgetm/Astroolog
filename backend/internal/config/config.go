package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func Load() (Config, error) {
	production := os.Getenv("APP_ENV") == "production"
	if !production {
		_ = godotenv.Load()
	}
	if production {
		for _, key := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"} {
			if os.Getenv(key) == "" {
				return Config{}, fmt.Errorf("required environment variable %s is missing", key)
			}
		}
		if os.Getenv("DB_SSLMODE") != "verify-full" {
			return Config{}, fmt.Errorf("production DB_SSLMODE must be verify-full")
		}
	}

	return Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5433"),
		DBUser:     getEnv("DB_USER", "astroolog"),
		DBPassword: getEnv("DB_PASSWORD", "astroolog"),
		DBName:     getEnv("DB_NAME", "astroolog"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
