package config

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	return &Config{
		//	os.Getenv(""),
	}
}
