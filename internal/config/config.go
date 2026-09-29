package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds every configurable value the application needs at startup.
type Config struct {
	ServerPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret      string
	JWTExpiryHours int
}

// Load reads a .env file (if present) and then environment variables,
// falling back to sensible defaults for local development.
func Load() Config {
	// It's fine if .env doesn't exist (e.g. in production where real
	// env vars are injected by the platform).
	_ = godotenv.Load()

	return Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "university_events"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiryHours: getEnvInt("JWT_EXPIRY_HOURS", 72),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	var result int
	if _, err := fmt.Sscanf(v, "%d", &result); err != nil {
		return fallback
	}
	return result
}
