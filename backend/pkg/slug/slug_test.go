package slug_test

import (
	"testing"

	"luny.dev/cherryauctions/pkg/slug"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple lowercase",
			input:    "hello world",
			expected: "hello-world",
		},
		{
			name:     "Uppercase and symbols",
			input:    "Hello, World!",
			expected: "hello-world",
		},
		{
			name:     "Numbers and multiple spaces",
			input:    "Item  Number  123",
			expected: "item-number-123",
		},
		{
			name:     "Trailing and leading special characters",
			input:    "---Service Name!!!",
			expected: "service-name",
		},
		{
			name:     "Complex characters (ASCII only)",
			input:    "Computers & Laptops @ 2025",
			expected: "computers-laptops-2025",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := slug.Slugify(tt.input)
			if actual != tt.expected {
				t.Errorf("Slugify(%q) = %q; want %q", tt.input, actual, tt.expected)
			}
		})
	}
}
