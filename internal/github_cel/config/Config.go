package config

import (
	"os"
)

type Config struct {
	Port string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	cfg.Port = getEnv("PORT", "50051")
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
