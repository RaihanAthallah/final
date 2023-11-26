package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EmailConfig struct {
	MailHost  string
	MailPort  int
	AuthEmail string
	AuthPass  string
}

func NewEmailConfig() (*EmailConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading .env file: %v\n", err)
	}

	emailConfig := &EmailConfig{
		MailHost:  os.Getenv("SMTP_HOST"),
		MailPort:  GetenvInt("SMTP_PORT"),
		AuthEmail: os.Getenv("SMTP_AUTH_EMAIL"),
		AuthPass:  os.Getenv("SMTP_AUTH_PASSWORD"),
	}

	return emailConfig, err
}
