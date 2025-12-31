package ranges_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"luny.dev/cherryauctions/pkg/ranges"
)

func TestEach(t *testing.T) {
	vals := []int{1, 2, 3}
	doubles := ranges.Each(vals, func(num int) int { return num * 2 })
	assert.ElementsMatch(t, doubles, []int{2, 4, 6})
}
