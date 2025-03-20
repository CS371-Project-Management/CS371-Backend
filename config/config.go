package config

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadConfig() error {
	return godotenv.Load()
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("Required environment variable not set: " + key)
	}
	return value
}
