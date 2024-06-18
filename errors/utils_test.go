package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_removeControlChars(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"regular string",
			"Hello world!",
			"Hello world!",
		},
		{
			"with control chars",
			"\t\tHello world!\n",
			"Hello world!",
		},

		{
			"control chars as text",
			"\\t\\tHello world!\\n",
			"\\t\\tHello world!\\n",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := removeControlChars(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_isAlphaSpace(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"lowercase string",
			"hello",
			true,
		},
		{
			"uppercase string",
			"HELLO",
			true,
		},
		{
			"first word capitalized",
			"Hello",
			true,
		},
		{
			"title case",
			"Hello World",
			true,
		},
		{
			"sentence case",
			"Hello world",
			true,
		},
		{
			"with control chars",
			"\t\tHello world!\n",
			false,
		},

		{
			"with special characters",
			"Hello!",
			false,
		},
		{
			"with numbers",
			"Hello123",
			false,
		},
		{
			"with spaces",
			"Hello World",
			true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := isAlphaSpace(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_isAlphaNum(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"lowercase string",
			"hello",
			true,
		},
		{
			"uppercase string",
			"HELLO",
			true,
		},
		{
			"first word capitalized",
			"Hello",
			true,
		},
		{
			"title case",
			"Hello World",
			false,
		},
		{
			"sentence case",
			"Hello world",
			false,
		},
		{
			"with control chars",
			"\t\tHello world!\n",
			false,
		},

		{
			"with special characters",
			"Hello!",
			false,
		},
		{
			"with numbers",
			"Hello123",
			true,
		},
		{
			"with spaces",
			"Hello World",
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := isAlphaNum(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_containsControl(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"without control chars",
			"hello",
			false,
		},
		{
			"with control chars",
			"\t\tHello world!\n",
			true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := containsControl(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_containsAlpha(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"with lowercase letter",
			"123a",
			true,
		},
		{
			"with uppercase letter",
			"123A",
			true,
		},
		{
			"without letter",
			"123",
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := containsAlpha(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_containsLower(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"with lowercase letter",
			"123a",
			true,
		},
		{
			"with uppercase letter",
			"123A",
			false,
		},
		{
			"without letter",
			"123",
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := containsLower(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_containsUpper(t *testing.T) {
	testcases := []struct {
		name     string
		input    byte
		expected bool
	}{
		{
			"lowercase letter",
			'a',
			false,
		},
		{
			"uppercase letter",
			'A',
			true,
		},
		{
			"without letter",
			'1',
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := isCharUpper(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_utilsisEmpty(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"empty string",
			"",
			true,
		},
		{
			"only spaces",
			"    ",
			true,
		},
		{
			"with one or more characters",
			"a1@",
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := isEmpty(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_isTitle(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"lowercase string",
			"hello",
			false,
		},
		{
			"uppercase string",
			"HELLO",
			false,
		},
		{
			"first word capitalized",
			"Hello",
			true,
		},
		{
			"title case",
			"Hello World",
			true,
		},
		{
			"sentence case",
			"Hello world",
			false,
		},
		{
			"with control chars",
			"\t\tHello world!\n",
			false,
		},

		{
			"with special characters",
			"Hello!",
			false,
		},
		{
			"with numbers",
			"Hello123",
			false,
		},
		{
			"with spaces",
			"Hello World",
			true,
		},
		{
			"leading spaces",
			"   Hello World",
			false,
		},
		{
			"trailing spaces",
			"Hello World   ",
			false,
		},
		{
			"consecutive spaces",
			"Hello      World",
			false,
		},
		{
			"empty string",
			"",
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := isTitle(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_isSentence(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"lowercase string",
			"hello",
			false,
		},
		{
			"uppercase string",
			"HELLO",
			false,
		},
		{
			"first word capitalized",
			"Hello",
			true,
		},
		{
			"title case",
			"Hello World",
			false,
		},
		{
			"sentence case",
			"Hello world",
			true,
		},
		{
			"with control chars",
			"\t\tHello world!\n",
			false,
		},

		{
			"with special characters",
			"Hello!",
			false,
		},
		{
			"with numbers",
			"Hello123",
			false,
		},
		{
			"with spaces",
			"Hello world",
			true,
		},
		{
			"leading spaces",
			"   Hello world",
			false,
		},
		{
			"trailing spaces",
			"Hello world   ",
			true,
		},
		{
			"consecutive spaces",
			"Hello      world",
			true,
		},
		{
			"empty string",
			"",
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := isSentence(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_hasConsecutiveSpaces(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"empty string",
			"",
			false,
		},
		{
			"with consecutive spaces",
			"    ",
			true,
		},
		{
			"letters and consecutive spaces",
			"a    b     c",
			true,
		},
		{
			"no consecutive spaces",
			"a b c",
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := hasConsecutiveSpaces(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}

func Test_hasTrailingOrLeading(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"empty string",
			"",
			false,
		},
		{
			"with consecutive spaces",
			"    ",
			true,
		},
		{
			"letters and consecutive spaces",
			"a    b     c",
			false,
		},
		{
			"letter with trailing spaces",
			"a b c   ",
			true,
		},
		{
			"letter with leading spaces",
			"    a b c",
			true,
		},
		{
			"without leading or trailing spaces",
			"a b c",
			false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result := hasTrailingOrLeading(tc.input)
			assert.Equal(result, tc.expected)
		})
	}
}
