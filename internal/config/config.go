package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Database struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

type ConfigApplication struct {
	Port string
}

// Config holds the application configuration
type Config struct {
	Database    Database
	Application ConfigApplication
}

// LoadConfig loads configuration from .env file
func LoadConfig() *Config {
	err := godotenv.Load("config.env")

	if err != nil {
		log.Fatalf("Error loading .env file%v", err)
	}

	return &Config{
		Database: Database{
			DBUser:     os.Getenv("DB_USER"),
			DBPassword: os.Getenv("DB_PASSWORD"),
			DBHost:     os.Getenv("DB_HOST"),
			DBPort:     os.Getenv("DB_PORT"),
			DBName:     os.Getenv("DB_NAME"),
		},
		Application: ConfigApplication{
			Port: os.Getenv("APP_PORT"),
		},
	}
}
