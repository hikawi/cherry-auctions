package closer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"luny.dev/cherryauctions/pkg/closer"
)

type TestCloser struct {
	closeCalled bool
}

func (c *TestCloser) Close() error {
	c.closeCalled = true
	return nil
}

func TestCloseResources(t *testing.T) {
	testCloser := &TestCloser{closeCalled: false}
	closer.CloseResources(testCloser)
	assert.Equal(t, testCloser.closeCalled, true)
}
