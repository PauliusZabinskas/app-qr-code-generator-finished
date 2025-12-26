package config

import (
	"log"
	"os"
	"strings"
)

// Config holds all application configuration
type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Security
	JWTSecret     string
	EncryptionKey string // Must be 32 bytes for AES-256

	// Server
	Port        string
	FrontendURL string
	GinMode     string
	LogLevel    string

	// CORS
	AllowedOrigins []string
}

// Load reads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		// Database defaults
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "wifiqr"),
		DBPassword: getEnv("DB_PASSWORD", "dev_password_123"),
		DBName:     getEnv("DB_NAME", "wifiqr_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// Security
		JWTSecret:     getEnv("JWT_SECRET", ""),
		EncryptionKey: getEnv("ENCRYPTION_KEY", ""),

		// Server defaults
		Port:        getEnv("PORT", "8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:4200"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),

		// CORS
		AllowedOrigins: parseAllowedOrigins(getEnv("ALLOWED_ORIGINS", "http://localhost:4200")),
	}

	// Validate required configuration
	cfg.validate()

	return cfg
}

// validate ensures all required configuration is present
func (c *Config) validate() {
	if c.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	if len(c.JWTSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters long")
	}

	if c.EncryptionKey == "" {
		log.Fatal("ENCRYPTION_KEY environment variable is required")
	}

	if len(c.EncryptionKey) != 32 {
		log.Fatal("ENCRYPTION_KEY must be exactly 32 characters (256 bits) for AES-256")
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// parseAllowedOrigins parses comma-separated origins
func parseAllowedOrigins(origins string) []string {
	if origins == "" {
		return []string{}
	}

	parts := strings.Split(origins, ",")
	result := make([]string, 0, len(parts))

	for _, origin := range parts {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
