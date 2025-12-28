package infra

import (
	"gopkg.in/gomail.v2"
)

func SetupMailer(smtpHost string, smtpPort int, smtpUser, smtpPassword string) *gomail.Dialer {
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)
	return dialer
}
