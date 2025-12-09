package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"luny.dev/cherryauctions/utils"
)

func TestGetenvDefault(t *testing.T) {
	val := utils.Getenv("GO_TEST_ENV", "hello")
	assert.Equal(t, "hello", val)
}

func TestGetenv(t *testing.T) {
	t.Setenv("GO_TEST_ENV", "world")
	val := utils.Getenv("GO_TEST_ENV", "hello")
	assert.Equal(t, "world", val)
}
