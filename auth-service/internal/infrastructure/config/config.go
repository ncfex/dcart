package config

import (
	"os"
)

type Config struct {
	PostgresURL string
	RedisURL    string
	Port        string
}

func LoadConfig() *Config {
	return &Config{
		PostgresURL: getEnv("POSTGRES_URL", "postgres://user:password@localhost:5432/authdb?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
		Port:        getEnv("AUTH_SERVICE_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
