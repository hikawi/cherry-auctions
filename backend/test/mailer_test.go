package test

import (
	"testing"

	"gopkg.in/gomail.v2"
	"luny.dev/cherryauctions/services"
)

func TestMailerService(t *testing.T) {
	dialer := services.NewMailerService()

	msg := gomail.NewMessage()
	msg.SetHeader("From", "Noreply <test@example.com>")
	msg.SetHeader("To", "Recipient <recipient@example.com>")
	msg.SetHeader("Subject", "Test Subject")
	msg.SetBody("text/plain", "Test Body")

	err := dialer.DialAndSend(msg)
	if err != nil {
		t.Fatalf("failed to send email: %v", err)
	}
}
