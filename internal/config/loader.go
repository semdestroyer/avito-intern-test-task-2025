package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	return &Config{
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_NAME:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_HOST"),
		DB_PASSWORD: os.Getenv("DB_HOST"),
		DB_USERNAME: os.Getenv("DB_HOST"),
	}
}
