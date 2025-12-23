// Package test provides some endpoints meant for testing only, and will not be available in production.
package test

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

type TestHandler struct {
	S3Client   *s3.Client
	MailDialer *gomail.Dialer
}

func (h *TestHandler) SetupRouter(g *gin.RouterGroup) {
	if gin.Mode() == gin.ReleaseMode {
		return
	}

	r := g.Group("/test")

	r.POST("", h.PostTest)
}
