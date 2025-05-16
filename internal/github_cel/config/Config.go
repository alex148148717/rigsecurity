package config

import (
	"os"
)

type Config struct {
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
