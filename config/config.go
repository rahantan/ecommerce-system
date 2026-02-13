package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Database struct {
	Driver   string
	DBName   string
	Host     string
	User     string
	Password string
	Port     string
}
type Server struct {
	Host string
	Port string
}
type Jwt struct {
	SecretKey string
}
type Config struct {
	Server
	Database
	Jwt
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("env file not found, using system env")
	}
	return &Config{
		Database: Database{
			Driver:   os.Getenv("DB_DRIVER"),
			Port:     os.Getenv("DB_PORT"),
			Host:     os.Getenv("DB_HOST"),
			DBName:   os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Jwt: Jwt{
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
		},
	}
}
