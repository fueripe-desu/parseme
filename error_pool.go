package parseme

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

func (p *ErrorPool) Error(data ErrorData) {
	p.errorStack.Push(data)
	p.Notify()
}

func (p *ErrorPool) AddError(data ErrorData) {
	p.errorStack.Push(data)
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
