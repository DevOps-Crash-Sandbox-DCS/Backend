package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv                string
	HTTPPort              string
	DatabaseURL           string
	JWTSecret             string
	JWTTTLHours           int
	MLHintsBaseURL        string
	MLHintsTimeoutSeconds int
}

func Load() Config {
	jwtTTLHours, err := strconv.Atoi(getEnv("JWT_TTL_HOURS", "24"))
	if err != nil {
		jwtTTLHours = 24
	}

	return Config{
		AppEnv:                getEnv("APP_ENV", "development"),
		HTTPPort:              getEnv("HTTP_PORT", "8080"),
		DatabaseURL:           getEnv("DATABASE_URL", ""),
		JWTSecret:             getEnv("JWT_SECRET", "dev_secret"),
		JWTTTLHours:           jwtTTLHours,
		MLHintsBaseURL:        getEnv("ML_HINTS_BASE_URL", "http://localhost:8000"),
		MLHintsTimeoutSeconds: getEnvInt("ML_HINTS_TIMEOUT_SECONDS", 10),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
