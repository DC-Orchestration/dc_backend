package config

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD") 
	smtpHost := os.Getenv("SMTP_HOST")   
	smtpPort := os.Getenv("SMTP_PORT")    

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("SMTP config missing from environment variables")
	}

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, password, smtpHost)

	addr := smtpHost + ":" + smtpPort
	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}
