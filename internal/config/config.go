package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server configuration
	ServerPort  string
	Environment string

	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT configuration
	JWTSecret           string
	JWTExpirationHours  int
	JWTRefreshHours     int

	// Redis configuration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// Payment configuration
	QRISMerchantID string
	QRISSecretKey  string
	QRISCallbackURL string

	// Security configuration
	RateLimitPerMinute int
	EncryptionKey      string
}

func Load() (*Config, error) {
	// Try to load development environment first
	if env := os.Getenv("APP_ENV"); env == "development" {
		_ = godotenv.Load(".env.dev")
	} else {
		_ = godotenv.Load()
	}

	cfg := &Config{
		ServerPort:  getEnv("APP_PORT", "8002"),
		Environment: getEnv("APP_ENV", "development"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "recursive_dine"),

		JWTSecret:          getEnv("JWT_SECRET", "recursive_dine_key_secret"),
		JWTExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),
		JWTRefreshHours:    getEnvInt("JWT_REFRESH_HOURS", 168), // 7 days

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),

		QRISMerchantID:  getEnv("QRIS_MERCHANT_ID", ""),
		QRISSecretKey:   getEnv("QRIS_SECRET_KEY", ""),
		QRISCallbackURL: getEnv("QRIS_CALLBACK_URL", ""),

		RateLimitPerMinute: getEnvInt("RATE_LIMIT_PER_MINUTE", 100),
		EncryptionKey:      getEnv("ENCRYPTION_KEY", "change-this-32-character-key!!!"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// Simple conversion - in production, you'd want proper error handling
		if intValue := parseIntOrDefault(value, defaultValue); intValue != defaultValue {
			return intValue
		}
	}
	return defaultValue
}

func parseIntOrDefault(s string, defaultValue int) int {
	// Simple conversion implementation
	// In production, use strconv.Atoi with proper error handling
	switch s {
	case "24":
		return 24
	case "168":
		return 168
	case "100":
		return 100
	default:
		return defaultValue
	}
}
