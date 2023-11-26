package utils

import (
	"a21hc3NpZ25tZW50/config"
	"log"

	gomail "gopkg.in/gomail.v2"
)

func SendEmail(dst string, subject string, body string) error {
	emailConfig, err := config.NewEmailConfig()

	if err != nil {
		log.Printf("error retrieving config: %v\n", err)
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", emailConfig.AuthEmail)
	msg.SetHeader("To", dst)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	n := gomail.NewDialer(
		emailConfig.MailHost, emailConfig.MailPort,
		emailConfig.AuthEmail, emailConfig.AuthPass)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		log.Printf("Failed to Send: %v\n", err)
		// log.Printf("%v %v \n", emailConfig.AuthEmail, emailConfig.AuthPass)
	}

	return nil
}
