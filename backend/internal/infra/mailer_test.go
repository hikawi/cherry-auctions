package infra_test

import (
	"testing"

	"gopkg.in/gomail.v2"
	"luny.dev/cherryauctions/internal/infra"
	"luny.dev/cherryauctions/pkg/env"
)

func TestMailerService(t *testing.T) {
	dialer := infra.SetupMailer(env.Fatalenv("SMTP_HOST"), int(env.FatalenvInt("SMTP_PORT")), env.Fatalenv("SMTP_USER"), env.Fatalenv("SMTP_PASSWORD"))

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
