package errors

import "errors"

type errorStack struct {
	values []errorData
}

func (s *errorStack) push(value errorData) {
	s.values = append(s.values, value)
}

func (s *errorStack) peek() errorData {
	if s.isEmpty() {
		panic(errors.New("Cannot peek empty stack."))
	}

	return s.values[s.size()-1]
}

func (s *errorStack) clear() {
	s.values = s.values[:0]
}

func (s *errorStack) isEmpty() bool {
	return len(s.values) == 0
}

func (s *errorStack) size() int {
	return len(s.values)
}
