package errors

import (
	"regexp"
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

func isAlphaSpace(input string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z ]+$`)
	return regex.MatchString(input)
}

func isAlphaNum(input string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return regex.MatchString(input)
}

func containsControl(input string) bool {
	for _, r := range input {
		if unicode.IsControl(r) {
			return true
		}
	}

	return false
}

func containsAlpha(input string) bool {
	regex := regexp.MustCompile(`[a-zA-Z]`)
	return regex.MatchString(input)
}

func containsLower(input string) bool {
	regex := regexp.MustCompile(`[a-z]`)
	return regex.MatchString(input)
}

func isCharUpper(input byte) bool {
	return input >= 'A' && input <= 'Z'
}

func isEmpty(input string) bool {
	return strings.ReplaceAll(input, " ", "") == ""
}

func isTitle(input string) bool {
	regex := regexp.MustCompile(`^[A-Z][a-z]*(?: [A-Z][a-z]*)*$`)
	return regex.MatchString(input)
}

func hasConsecutiveSpaces(input string) bool {
	// Note the space before the curly braces, it indicates
	// the space character
	regex := regexp.MustCompile(` {2,}`)
	return regex.MatchString(input)
}

func isSentence(input string) bool {
	regex := regexp.MustCompile(`^[A-Z][a-z ]*$`)
	return regex.MatchString(input)
}

func hasTrailingOrLeading(input string) bool {
	if input == "" {
		return false
	}

	return input[0] == ' ' || input[len(input)-1] == ' '
}
