package errors

type errorStack struct {
	values []errorData
}

func (s *errorStack) Push(value errorData) {
	s.values = append(s.values, value)
}

func (s *errorStack) Pop() errorData {
	if s.IsEmpty() {
		panic("Cannot pop if stack is empty.")
	}

	popped := s.Peek()
	s.values = s.values[:s.Size()-1]

	return popped
}

func (s *errorStack) Peek() errorData {
	if s.IsEmpty() {
		panic("Cannot peek empty stack.")
	}

	return s.values[s.Size()-1]
}

func (s *errorStack) Clear() {
	s.values = s.values[:0]
}

func (s *errorStack) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *errorStack) Size() int {
	return len(s.values)
}
