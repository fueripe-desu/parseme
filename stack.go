package parseme

type stackData interface {
	preToken | Element | *Element | ErrorData
}

type stack[T stackData] struct {
	values []T
}

func (s *stack[T]) Push(value T) {
	s.values = append(s.values, value)
}

func (s *stack[T]) Pop() T {
	if s.IsEmpty() {
		panic(&StackError{message: "Cannot pop if stack is empty."})
	}

	popped := s.Peek()
	s.values = s.values[:s.Size()-1]

	return popped
}

func (s *stack[T]) Peek() T {
	if s.IsEmpty() {
		panic(&StackError{message: "Cannot peek empty stack."})
	}

	return s.values[s.Size()-1]
}

func (s *stack[T]) Clear() {
	s.values = s.values[:0]
}

func (s *stack[T]) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *stack[T]) Size() int {
	return len(s.values)
}
