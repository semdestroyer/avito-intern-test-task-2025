package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../.env")
		if err != nil {
			err = godotenv.Load("../../.env")
			if err != nil {
				log.Println("No .env file found, using default values or environment variables")
			}
		}
	}

	return &Config{
		DB_HOST:     getEnvOrDefault("DB_HOST", "localhost"),
		DB_NAME:     getEnvOrDefault("DB_NAME", "postgres"),
		DB_PORT:     getEnvOrDefault("DB_PORT", "5432"),
		DB_PASSWORD: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DB_USER:     getEnvOrDefault("DB_USER", "postgres"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
