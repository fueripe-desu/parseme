package parseme

import (
	"strconv"
)

type ErrorLevel int

const (
	Fatal ErrorLevel = iota
	Warning
	Info
	Debug
	Trace
)

type ErrorData struct {
	name    string
	level   ErrorLevel
	message string
}

type ErrorPool struct {
	observers  []ErrorObserver
	errorStack stack[ErrorData]
}

type ErrorObserver interface {
	OnUpdate(string, ErrorLevel, string)
}

func (p *ErrorPool) Subscribe(observer ErrorObserver) error {
	for _, o := range p.observers {
		if o == observer {
			return &ObserverDuplicateError{}
		}
	}

	p.observers = append(p.observers, observer)
	return nil
}

func (p *ErrorPool) Unsubscribe(observer ErrorObserver) error {
	newObservers := []ErrorObserver{}
	found := false

	for _, o := range p.observers {
		if o == observer {
			found = true
			continue
		}

		newObservers = append(newObservers, o)
	}

	if !found {
		return &ObserverNotFoundError{}
	}

	p.observers = newObservers
	return nil
}

func (p *ErrorPool) UnsubscribeAll() {
	p.observers = p.observers[:0]
}

func (p *ErrorPool) Error(data ErrorData, args []string) {
	p.AddError(data, args)
	p.Notify()
}

func (p *ErrorPool) AddError(data ErrorData, args []string) {
	var finalData ErrorData
	if args != nil && len(args) > 0 {
		finalData = p.precompileError(data, args)
	} else {
		finalData = data
	}

	p.errorStack.Push(finalData)
}

func (p *ErrorPool) precompileError(data ErrorData, args []string) ErrorData {
	newMsg := p.replaceMasks(data.message, args)
	return ErrorData{name: data.name, level: data.level, message: newMsg}
}

func (p *ErrorPool) replaceMasks(input string, args []string) string {
	chars := []byte(input)
	buffer := []byte{}

	for i := 0; i < len(chars); i++ {
		current := chars[i]

		if next := i + 1; current == '%' && next < len(chars) && chars[next] == '{' {
			newIndex, argIndex := p.consumeMask(i+2, &chars)
			i = newIndex

			if argIndex == -1 {
				continue
			}

			if argIndex >= len(args) {
				continue
			}

			argBytes := []byte(args[argIndex])
			buffer = append(buffer, argBytes...)
		} else {
			buffer = append(buffer, current)
		}
	}

	return string(buffer)
}

func (p *ErrorPool) consumeMask(start int, chars *[]byte) (int, int) {
	buffer := []byte{}

	startIndex := start
	for ; startIndex < len(*chars); startIndex++ {
		current := (*chars)[startIndex]

		if current == '}' {
			break
		}

		if current >= '0' && current <= '9' {
			buffer = append(buffer, (*chars)[startIndex])
		}
	}

	var argsIndex int

	if len(buffer) == 0 {
		argsIndex = -1
	} else {
		intIndex, _ := strconv.Atoi(string(buffer))
		argsIndex = intIndex
	}

	return startIndex, argsIndex
}

func (p *ErrorPool) Notify() error {
	if !p.HasErrors() {
		return &NotifyObserverError{}
	}

	last := p.errorStack.Peek()

	for _, o := range p.observers {
		o.OnUpdate(last.name, last.level, last.message)
	}

	return nil
}

func (p *ErrorPool) ClearErrors() {
	p.errorStack.Clear()
}

func (p *ErrorPool) HasErrors() bool {
	return !p.errorStack.IsEmpty()
}
