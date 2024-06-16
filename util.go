package parseme

import (
	"regexp"
	"strings"
	"unicode"
)

func isValidIdentifier(value string) bool {
	isValid := hasValidIdentifierChars(value)
	firstChar := firstCharIsDigit(value)
	return isValid && (!firstChar)
}

func isValidProperty(value string) bool {
	isValid := hasValidPropertyChars(value)
	firstChar := firstCharIsDigit(value)
	return isValid && (!firstChar)
}

func hasValidPropertyChars(value string) bool {
	if len(value) == 0 {
		return false
	}

	pattern := `^[a-zA-Z0-9\-\_]*$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(value)
}

func hasValidIdentifierChars(value string) bool {
	if len(value) == 0 {
		return false
	}

	pattern := `^[a-zA-Z0-9\-\:\.\_]*$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(value)
}

func firstCharIsDigit(value string) bool {
	if len(value) == 0 {
		return false
	}

	pattern := `[0-9]`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(string(value[0]))
}

func isDoubleQuoted(value string) bool {
	if len(value) < 2 {
		return false
	}

	return value[0] == '"' && value[len(value)-1] == '"'
}

func isSingleQuoted(value string) bool {
	if len(value) < 2 {
		return false
	}

	return value[0] == '\'' && value[len(value)-1] == '\''
}

func removeQuotes(value string) string {
	if isSingleQuoted(value) || isDoubleQuoted(value) {
		return value[1 : len(value)-1]
	}

	return value
}

func removeControlChars(input string) string {
	var builder strings.Builder
	for _, r := range input {
		if !unicode.IsControl(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
