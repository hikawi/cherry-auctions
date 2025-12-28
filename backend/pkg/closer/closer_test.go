package closer_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"luny.dev/cherryauctions/pkg/closer"
)

type TestCloser struct {
	makeFail    bool
	closeCalled bool
}

func (c *TestCloser) Close() error {
	if c.makeFail {
		return errors.New("lol")
	}

	c.closeCalled = true
	return nil
}

func TestCloseResources(t *testing.T) {
	testCloser := &TestCloser{closeCalled: false}
	closer.CloseResources(testCloser)
	assert.Equal(t, testCloser.closeCalled, true)
}
