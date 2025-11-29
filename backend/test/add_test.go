package test

import "testing"

func TestAddition(t *testing.T) {
	val := 1 + 1
	if val != 2 {
		t.Errorf("1 + 1 = %d; want 2", val)
	}
}
