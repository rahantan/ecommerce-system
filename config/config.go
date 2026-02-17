package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

type Midtrans struct {
	ServerKey string
	Env       string
	snap.Client
}
type Config struct {
	Server
	Database
	Jwt
	Midtrans
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
		Midtrans: Midtrans{
			ServerKey: os.Getenv("MIDTRANS_SERVER_KEY"),
			Env:       os.Getenv("MIDTRANS_ENV"),
		},
	}
}
func (conf *Config) ConnectionDb() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return db
}

func (conf *Config) InitMidtrans() {
	env := midtrans.Sandbox

	if conf.Midtrans.Env == "production" {
		env = midtrans.Production
	}

	conf.Midtrans.Client.New(conf.Midtrans.ServerKey, env)
}
