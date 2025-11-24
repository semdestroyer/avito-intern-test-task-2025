package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadConfig() *Config {
	if err := godotenv.Load("../../.env"); err != nil { //TODO: сделать более красиво
		log.Fatal("No .env file found:", err)
	}

	return &Config{
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_USER:     os.Getenv("DB_USER"),
	}
}
