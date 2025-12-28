package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"luny.dev/cherryauctions/internal/config"
	"luny.dev/cherryauctions/internal/repositories"
	"luny.dev/cherryauctions/internal/routes"
	"luny.dev/cherryauctions/internal/services"
)

func TestHealth(t *testing.T) {
	server := gin.New()
	routes.SetupRoutes(server, routes.ServerDependency{
		Version:      "v1",
		DB:           nil,
		S3Client:     nil,
		MailDialer:   nil,
		Config:       &config.Config{},
		Services:     services.ServiceRegistry{},
		Repositories: repositories.RepositoryRegistry{},
	})
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/v1/health", nil)
	server.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"healthy\"}", w.Body.String())
}
