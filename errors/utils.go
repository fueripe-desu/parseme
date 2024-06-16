package errors

import (
	"strings"
	"unicode"
)

func removeControlChars(input string) string {
	var builder strings.Builder
	for _, r := range input {
		if !unicode.IsControl(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
