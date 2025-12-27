package utils

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`[^a-z0-9]+`)

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = re.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
