package config

import (
	"os"
	"strconv"
	"time"
)

// AppConfig holds the application-wide settings.
type AppConfig struct {
	DBPath          string
	APIPort         string
	SMTPSender      string
	RateLimitPerMin      int
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// LoadConfig loads configuration from environment variables with defaults.
func LoadConfig() *AppConfig {
	return &AppConfig{
		DBPath:               getEnv("DB_PATH", "emails.db"),
		APIPort:              getEnv("API_PORT", "8080"),
		SMTPSender:           getEnv("SMTP_SENDER", "sanjanamaahi2001@gmail.com"),
		RateLimitPerMin:      getEnvInt("RATE_LIMIT_IP_MIN", 5),
		JWTSecret:            getEnv("JWT_SECRET", "change-me-in-production"),
		AccessTokenDuration:  getEnvDuration("ACCESS_TOKEN_DURATION", 15*time.Minute),
		RefreshTokenDuration: getEnvDuration("REFRESH_TOKEN_DURATION", 168*time.Hour), // 7 days
	}
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return fallback
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
