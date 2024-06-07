package parseme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isValidIdentifier(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"valid identifier",
			"valid",
			true,
		},
		{
			"starts with digit",
			"1nvalid",
			false,
		},
		{
			"empty value",
			"",
			false,
		},
		{
			"numbers after first char",
			"a123",
			true,
		},
		{
			"contains invalid char",
			"asd@",
			false,
		},
		{
			"contains hyphen",
			"-b-c",
			true,
		},
		{
			"contains dot",
			".b.c",
			true,
		},
		{
			"contains colon",
			":b:c",
			true,
		},

		{
			"contains underscore",
			"_b_c",
			true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := isValidIdentifier(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_isValidProperty(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"valid identifier",
			"valid",
			true,
		},
		{
			"starts with digit",
			"1nvalid",
			false,
		},
		{
			"empty value",
			"",
			false,
		},
		{
			"numbers after first char",
			"a123",
			true,
		},
		{
			"contains invalid char",
			"asd@",
			false,
		},
		{
			"contains hyphen",
			"-b-c",
			true,
		},
		{
			"contains dot",
			".b.c",
			false,
		},
		{
			"contains colon",
			":b:c",
			false,
		},

		{
			"contains underscore",
			"_b_c",
			true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := isValidProperty(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_hasValidIdentifierChars(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"contains valid chars",
			"valid",
			true,
		},
		{
			"contains invalid char",
			"asd@",
			false,
		},

		{
			"empty value",
			"",
			false,
		},
		{
			"contains hyphen",
			"-b-c",
			true,
		},
		{
			"contains dot",
			".b.c",
			true,
		},
		{
			"contains colon",
			":b:c",
			true,
		},

		{
			"contains underscore",
			"_b_c",
			true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := hasValidIdentifierChars(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_hasValidPropertyChars(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"contains valid chars",
			"valid",
			true,
		},
		{
			"contains invalid char",
			"asd@",
			false,
		},

		{
			"empty value",
			"",
			false,
		},
		{
			"contains hyphen",
			"-b-c",
			true,
		},
		{
			"contains underscore",
			"_b_c",
			true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := hasValidPropertyChars(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_firstCharIsDigit(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"first char is digit",
			"1valid",
			true,
		},
		{
			"first char is not digit",
			"asd",
			false,
		},
		{
			"empty value",
			"",
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := firstCharIsDigit(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_isDoubleQuoted(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"double quoted value",
			"\"something\"",
			true,
		},
		{
			"single quoted value",
			"'something'",
			false,
		},
		{
			"unquoted value",
			"something",
			false,
		},
		{
			"single double quotes",
			"'something\"",
			false,
		},
		{
			"double single quotes",
			"\"something'",
			false,
		},
		{
			"empty value",
			"",
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := isDoubleQuoted(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_isSingleQuoted(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"double quoted value",
			"\"something\"",
			false,
		},
		{
			"single quoted value",
			"'something'",
			true,
		},
		{
			"unquoted value",
			"something",
			false,
		},
		{
			"single double quotes",
			"'something\"",
			false,
		},
		{
			"double single quotes",
			"\"something'",
			false,
		},
		{
			"empty value",
			"",
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := isSingleQuoted(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_removeQuotes(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"double quoted value",
			"\"something\"",
			"something",
		},
		{
			"single quoted value",
			"'something'",
			"something",
		},
		{
			"unquoted value",
			"something",
			"something",
		},
		{
			"single double quotes",
			"'something\"",
			"'something\"",
		},
		{
			"double single quotes",
			"\"something'",
			"\"something'",
		},
		{
			"empty value",
			"",
			"",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := removeQuotes(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}
