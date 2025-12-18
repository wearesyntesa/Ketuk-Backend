package utils

import (
	"net/smtp"
)

func SendEmail(to []string, subject, body string, auth smtp.Auth, host string, from string) error{
	msg := "From: " + from + "\n" +
		"To: " + to[0] + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail(host+":587", auth, from, to, []byte(msg))
	return err
}

func SetSMTPAuth(email, password, host string) smtp.Auth {
	return smtp.PlainAuth("", email, password, host)
}
