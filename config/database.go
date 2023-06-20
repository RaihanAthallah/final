package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func GetenvInt(envName string) int {
	strEnv := os.Getenv(envName)
	i, err := strconv.Atoi(strEnv)
	if err != nil {
		log.Printf("error parsing env %v to int: %v\n", envName, err)
	}
	return i
}

type config struct {
	AppPort    string
	DBHost     string
	DBUsername string
	DBPassword string
	DBName     string
	DBPort     int
	JWTKey     []byte
}

var Config *config

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading .env file: %v\n", err)
	}

	if Config == nil {
		Config = &config{
			AppPort:    os.Getenv("PORT"),
			DBHost:     os.Getenv("DB_HOST"),
			DBUsername: os.Getenv("DB_USERNAME"),
			DBPassword: os.Getenv("DB_PASSWORD"),
			DBName:     os.Getenv("DB_NAME"),
			DBPort:     GetenvInt("DB_PORT"),
			JWTKey:     []byte(os.Getenv("JWT_SECRET")),
		}
	}
}
