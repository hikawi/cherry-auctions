package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"luny.dev/cherryauctions/routes"
)

func TestHealth(t *testing.T) {
	server := gin.Default()
	routes.SetupRoutes(server, routes.ServerDependency{
		DB:         nil,
		S3Client:   nil,
		MailDialer: nil,
	})
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/v1/health", nil)
	server.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"healthy\"}", w.Body.String())
}
