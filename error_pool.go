package parseme

import (
	"bufio"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type ErrorLevel int

const (
	Fatal ErrorLevel = iota
	Error
	Warning
	Info
	Debug
	Trace
)

type errorData struct {
	name    string
	message string
	code    string
	module  string
	fix     string
	level   ErrorLevel
}

type ErrorInfo struct {
	Name    string
	Message string
	Code    string
	Module  string
	Fix     string

	CallerFuncName string
	CallerFilename string
	CallerLine     int

	ErrorFuncName string
	ErrorFilename string
	ErrorLine     int
	ErrorSite     string

	Level     ErrorLevel
	Timestamp time.Time
}

type ErrorPool struct {
	observers  []ErrorObserver
	errorStack stack[errorData]
}

type ErrorObserver interface {
	OnUpdate(info ErrorInfo)
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

func (p *ErrorPool) Error(data *errorData, args []string) {
	if data == nil {
		p.Error(nilErrorDataError, nil)
	} else {
		p.AddError(*data, args)
		p.Notify()
	}
}

func (p *ErrorPool) AddError(data errorData, args []string) {
	var finalData errorData
	if args != nil && len(args) > 0 {
		finalData = p.precompileError(data, args)
	} else {
		finalData = data
	}

	p.errorStack.Push(finalData)
}

func (p *ErrorPool) precompileError(data errorData, args []string) errorData {
	newMsg := p.replaceMasks(data.message, args)
	return errorData{name: data.name, level: data.level, message: newMsg}
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

	recursive := false
	if last.name == "Nil error data" {
		recursive = true
	}

	callerLine, callerFile, callerFuncName, _ := p.getErrorInfo(recursive, true)
	errorLine, errorFile, errorFuncName, errorSite := p.getErrorInfo(recursive, false)

	for _, o := range p.observers {
		info := ErrorInfo{
			Name:           last.name,
			Message:        last.message,
			Code:           last.code,
			Module:         last.module,
			Fix:            last.fix,
			CallerFuncName: callerFuncName,
			CallerFilename: callerFile,
			CallerLine:     callerLine,
			ErrorFuncName:  errorFuncName,
			ErrorFilename:  errorFile,
			ErrorLine:      errorLine,
			ErrorSite:      errorSite,
			Level:          last.level,
			Timestamp:      time.Now(),
		}
		o.OnUpdate(info)
	}

	return nil
}

func (p *ErrorPool) getErrorInfo(recursive bool, caller bool) (int, string, string, string) {
	index := p.getStackIndex(recursive, caller)

	pc, file, line, ok := runtime.Caller(index)
	funcName := p.extractFuncName(runtime.FuncForPC(pc).Name())

	if !ok {
		panic("Could not retrieve runtime caller.")
	}

	return line, filepath.Base(file), funcName, p.getLineContents(file, line)
}

func (p *ErrorPool) getStackIndex(recursive bool, caller bool) int {
	if recursive {
		if caller {
			return 3
		}

		return 4
	} else {
		if caller {
			return 2
		}

		return 3
	}
}

func (p *ErrorPool) extractFuncName(fullName string) string {
	splitted := strings.Split(fullName, "/")
	last := splitted[len(splitted)-1]
	pos := strings.Index(last, ".")
	return last[pos+1:]
}

func (p *ErrorPool) getLineContents(filepath string, line int) string {
	file, err := os.Open(filepath)

	if err != nil {
		panic("Could not open file.")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 1

	for scanner.Scan() {
		if currentLine == line {
			withoutControl := removeControlChars(scanner.Text())
			return strings.Trim(withoutControl, " ")
		}

		currentLine++
	}

	if scannerErr := scanner.Err(); scannerErr != nil {
		panic("Scanner could not read the file.")
	}

	panic("File line does not exist.")
}

func (p *ErrorPool) ClearErrors() {
	p.errorStack.Clear()
}

func (p *ErrorPool) HasErrors() bool {
	return !p.errorStack.IsEmpty()
}

func NewErrorData(
	name string,
	message string,
	code string,
	module string,
	fix string,
	level ErrorLevel,
) *errorData {
	return &errorData{
		name:    name,
		message: message,
		code:    code,
		module:  module,
		fix:     fix,
		level:   level,
	}
}

func NewFatal(
	name string,
	message string,
	code string,
	module string,
	fix string,
) *errorData {
	return NewErrorData(
		name,
		message,
		code,
		module,
		fix,
		Fatal,
	)
}

func NewError(
	name string,
	message string,
	code string,
	module string,
	fix string,
) *errorData {
	return NewErrorData(
		name,
		message,
		code,
		module,
		fix,
		Error,
	)
}

func NewWarning(
	name string,
	message string,
	code string,
	module string,
	fix string,
) *errorData {
	return NewErrorData(
		name,
		message,
		code,
		module,
		fix,
		Warning,
	)
}

func NewInfo(
	message string,
	module string,
) *errorData {
	return NewErrorData(
		"",
		message,
		"",
		module,
		"",
		Info,
	)
}

func NewDebug(
	message string,
	module string,
) *errorData {
	return NewErrorData(
		"",
		message,
		"",
		module,
		"",
		Debug,
	)
}

func NewTrace(
	message string,
	module string,
) *errorData {
	return NewErrorData(
		"",
		message,
		"",
		module,
		"",
		Trace,
	)
}
