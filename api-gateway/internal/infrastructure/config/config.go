package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiGatewayPort string
	AuthServiceURL string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return &Config{}, err
	}
	return &Config{
		ApiGatewayPort: getEnv("API_GATEWAY_PORT", "8080"),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://auth-service:8080"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
