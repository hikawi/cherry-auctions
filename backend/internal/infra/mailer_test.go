package infra_test

import (
	"testing"

	"gopkg.in/gomail.v2"
	"luny.dev/cherryauctions/internal/config"
	"luny.dev/cherryauctions/internal/infra"
)

func TestMailerService(t *testing.T) {
	cfg := config.Load()
	dialer := infra.SetupMailer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.User, cfg.SMTP.Password)

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
