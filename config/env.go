package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JwtSecretKey string

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	JwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	if JwtSecretKey == "" {
		log.Fatal("JWT_SECRET_KEY is not set")
	}
}
