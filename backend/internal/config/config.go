package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AdminAllowedOrigins []string
	AppPort             string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBSSLMode           string
}

func Load() (Config, error) {
	production := os.Getenv("APP_ENV") == "production"
	if !production {
		_ = godotenv.Load()
	}
	if production {
		origin, err := url.Parse(os.Getenv("APP_ORIGIN"))
		if err != nil || origin.Scheme != "https" || origin.Host == "" || origin.Path != "" || origin.User != nil || origin.RawQuery != "" || origin.Fragment != "" {
			return Config{}, fmt.Errorf("production APP_ORIGIN must be an HTTPS origin")
		}
		for _, key := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"} {
			if os.Getenv(key) == "" {
				return Config{}, fmt.Errorf("required environment variable %s is missing", key)
			}
		}
		if os.Getenv("DB_SSLMODE") != "verify-full" {
			return Config{}, fmt.Errorf("production DB_SSLMODE must be verify-full")
		}
	}

	origins, err := adminOrigins(os.Getenv("APP_ENV"), os.Getenv("ADMIN_ALLOWED_ORIGINS"), os.Getenv("APP_ORIGIN"))
	if err != nil {
		return Config{}, err
	}
	return Config{
		AdminAllowedOrigins: origins,
		AppPort:             getEnv("APP_PORT", "8080"),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5433"),
		DBUser:              getEnv("DB_USER", "astroolog"),
		DBPassword:          getEnv("DB_PASSWORD", "astroolog"),
		DBName:              getEnv("DB_NAME", "astroolog"),
		DBSSLMode:           getEnv("DB_SSLMODE", "disable"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

// ValidOrigin rejects URLs that contain anything beyond scheme and authority.
func ValidOrigin(value string) bool {
	u, err := url.Parse(value)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Hostname() != "" && u.User == nil && u.Path == "" && u.RawQuery == "" && !u.ForceQuery && u.Fragment == "" && !strings.ContainsAny(value, "* \t\r\n#") && u.String() == value
}
func adminOrigins(environment, list, appOrigin string) ([]string, error) {
	if strings.TrimSpace(list) == "" {
		if appOrigin != "" {
			list = appOrigin
		} else if environment != "production" {
			list = "http://localhost:5174,http://127.0.0.1:5174"
		}
	}
	origins := []string{}
	for _, part := range strings.Split(list, ",") {
		origin := strings.TrimSpace(part)
		if origin == "" {
			continue
		}
		if !ValidOrigin(origin) || (environment == "production" && !strings.HasPrefix(origin, "https://")) {
			return nil, fmt.Errorf("ADMIN_ALLOWED_ORIGINS must contain exact HTTP(S) origins (HTTPS in production)")
		}
		origins = append(origins, origin)
	}
	return origins, nil
}
