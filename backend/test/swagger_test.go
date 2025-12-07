package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"luny.dev/cherryauctions/routes"
)

func TestSwaggerRedirect(t *testing.T) {
	server := gin.Default()
	routes.SetupRoutes(server, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/swagger", nil)
	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMovedPermanently, w.Code)
}
