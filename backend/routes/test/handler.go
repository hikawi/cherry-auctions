package test

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"luny.dev/cherryauctions/utils"
)

func (h *TestHandler) PostTest(g *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required,email"`
	}
	err := g.ShouldBindBodyWithJSON(&body)
	if err != nil {
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		utils.Log(gin.H{"message": err.Error()})
		return
	}

	msg := gomail.NewMessage()
	msg.SetHeaders(map[string][]string{
		"From": {fmt.Sprintf("CherryAuctions NoReply <%s>", utils.Fatalenv("SMTP_USER"))},
		"To":   {body.Email},
	})
	msg.SetBody("text/html", "<p>Hello there, <strong>mister handsome</strong></p>")
	err = h.MailDialer.DialAndSend(msg)

	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "couldn't send email"})
		utils.Log(gin.H{"message": err.Error()})
		return
	}

	g.JSON(http.StatusOK, gin.H{"message": "successfully sent email"})
}
