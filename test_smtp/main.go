package main

import (
	"log"
	"net/smtp"
)

func main() {
	// Set up authentication information.
	auth := smtp.PlainAuth("", "blueshoko@gmail.com", "lknj mkat xpcr jzkm", "smtp.gmail.com")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{"naufalfarras.dev@gmail.com"}
	msg := []byte("To: naufalfarras.dev@gmail.com\r\n" +
		"Subject: discount Gophers!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "blueshoko@gmail.com", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}
