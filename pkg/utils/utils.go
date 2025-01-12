package utils

import (
	"regexp"
	"strings"
)

// Polyfill
func HasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func CutSuffix(s, suffix string) (before string, found bool) {
	if !HasSuffix(s, suffix) {
		return s, false
	}
	return s[:len(s)-len(suffix)], true
}

func SlugifyString(input string) string {
	r := regexp.MustCompile(`(\(.*\))|(\[.*\])|(\.\w*$)|[^a-z0-9A-Z]`)
	rep := r.ReplaceAllStringFunc(input, func(m string) string {
		return ""
	})
	return strings.ToLower(rep)
}
