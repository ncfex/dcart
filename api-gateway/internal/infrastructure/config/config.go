package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig
	Services []ServiceConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type ServiceConfig struct {
	Name         string        `yaml:"name"`
	BaseURL      string        `yaml:"baseURL"`
	Timeout      time.Duration `yaml:"timeout"`
	RequiresAuth bool          `yaml:"requiresAuth"`
}

type AuthConfig struct {
	ServiceURL string        `yaml:"serviceURL"`
	Timeout    time.Duration `yaml:"timeout"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
