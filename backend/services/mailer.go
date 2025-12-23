package services

import (
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
	"luny.dev/cherryauctions/utils"
)

func NewMailerService() *gomail.Dialer {
	port, err := strconv.ParseInt(utils.Getenv("SMTP_PORT", "465"), 10, 0)
	if err != nil {
		log.Fatal(err)
	}

	dialer := gomail.NewDialer(utils.Fatalenv("SMTP_HOST"), int(port), utils.Fatalenv("SMTP_USER"), utils.Fatalenv("SMTP_PASSWORD"))
	return dialer
}
