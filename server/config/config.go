package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	AppEnv      string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://garden:garden@localhost:5432/garden?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
		AppEnv:      getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
