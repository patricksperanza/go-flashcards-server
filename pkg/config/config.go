package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JWTSecretKey []byte
var DBName string
var DBUsername string
var DBPassword string
var DBHost string
var DBPort string

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file. Using environment variables.")
	}
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("JWT_SECRET_KEY is not set")
	}
	JWTSecretKey = []byte(secret)

	DBName = os.Getenv("DB_NAME")
	DBUsername = os.Getenv("DB_USERNAME")
	DBPassword= os.Getenv("DB_PASSWORD")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	if DBUsername == "" || DBPassword == "" || DBHost == "" || DBPort == ""  || DBName == "" {
        log.Fatal("Missing one or more required environment variables")
    }
}
