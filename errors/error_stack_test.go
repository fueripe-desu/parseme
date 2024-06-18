package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_push(t *testing.T) {
	t.Run("add one element", func(t *testing.T) {
		assert := assert.New(t)
		stack := errorStack{}
		data := errorData{
			name:    "Not found error",
			message: "Could not find object.",
			code:    "ABC",
			fix:     "Try searching for the object.",
		}
		stack.push(data)
		assert.Equal(stack.values, []errorData{data})
	})

	t.Run("add more than one element", func(t *testing.T) {
		assert := assert.New(t)
		stack := errorStack{}
		data := errorData{
			name:    "Not found error",
			message: "Could not find object.",
			code:    "ABC",
			fix:     "Try searching for the object.",
		}
		stack.push(data)
		stack.push(data)
		stack.push(data)

		expected := []errorData{data, data, data}

		assert.Equal(stack.values, expected)
	})
}

func Test_peek(t *testing.T) {
	t.Run("peek last element", func(t *testing.T) {
		assert := assert.New(t)
		stack := errorStack{}
		data := errorData{
			name:    "Not found error",
			message: "Could not find object.",
			code:    "ABC",
			fix:     "Try searching for the object.",
		}
		stack.push(data)
		result := stack.peek()
		assert.Equal(result, data)
	})

	t.Run("peek empty stack", func(t *testing.T) {
		assert := assert.New(t)
		stack := errorStack{}
		assert.PanicsWithError("Cannot peek empty stack.", func() {
			stack.peek()
		})
	})
}

func Test_isEmpty(t *testing.T) {
	t.Run("non empty stack", func(t *testing.T) {
		assert := assert.New(t)
		stack := errorStack{}
		data := errorData{
			name:    "Not found error",
			message: "Could not find object.",
			code:    "ABC",
			fix:     "Try searching for the object.",
		}
		stack.push(data)
		result := stack.isEmpty()
		assert.Equal(result, false)
	})

	t.Run("empty stack", func(t *testing.T) {
		assert := assert.New(t)
		stack := errorStack{}
		result := stack.isEmpty()
		assert.Equal(result, true)
	})
}

func Test_size(t *testing.T) {
	t.Run("non empty stack", func(t *testing.T) {
		assert := assert.New(t)
		stack := errorStack{}
		data := errorData{
			name:    "Not found error",
			message: "Could not find object.",
			code:    "ABC",
			fix:     "Try searching for the object.",
		}
		stack.push(data)
		stack.push(data)
		stack.push(data)

		result := stack.size()
		assert.Equal(result, 3)
	})

	t.Run("empty stack", func(t *testing.T) {
		assert := assert.New(t)
		stack := errorStack{}
		result := stack.size()
		assert.Equal(result, 0)
	})
}

func Test_clear(t *testing.T) {
	t.Run("non empty stack", func(t *testing.T) {
		assert := assert.New(t)
		stack := errorStack{}
		data := errorData{
			name:    "Not found error",
			message: "Could not find object.",
			code:    "ABC",
			fix:     "Try searching for the object.",
		}
		stack.push(data)
		stack.push(data)
		stack.push(data)

		assert.Equal(stack.size(), 3)
		stack.clear()
		assert.Equal(stack.size(), 0)
	})

	t.Run("clear empty stack", func(t *testing.T) {
		assert := assert.New(t)
		stack := errorStack{}
		stack.clear()
		assert.Equal(stack.size(), 0)
	})
}
