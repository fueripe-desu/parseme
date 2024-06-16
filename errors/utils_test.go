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
