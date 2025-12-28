package env_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"luny.dev/cherryauctions/pkg/env"
)

func TestGetenvDefault(t *testing.T) {
	val := env.Getenv("GO_TEST_ENV", "hello")
	assert.Equal(t, "hello", val)
}

func TestGetenv(t *testing.T) {
	t.Setenv("GO_TEST_ENV", "world")
	val := env.Getenv("GO_TEST_ENV", "hello")
	assert.Equal(t, "world", val)
}

func TestGetenvInt(t *testing.T) {
	t.Setenv("GO_TEST_ENV", "1234")
	val := env.FatalenvInt("GO_TEST_ENV")
	assert.Equal(t, val, 1234)
}

func TestGetenvBool(t *testing.T) {
	t.Setenv("GO_TEST_ENV", "false")
	val := env.FatalenvInt("GO_TEST_ENV")
	assert.Equal(t, val, false)
}
