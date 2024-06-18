package errors

import (
	"bufio"
	"errors"
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
	errorStack errorStack
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

func (p *ErrorPool) error(level ErrorLevel, data *errorData, args []string) {
	if data == nil {
		panic(errors.New("Error data must not be nil."))
	}

	p.addError(data, args)
	p.notify(level)

	if level == Fatal {
		os.Exit(-1)
	}
}

func (p *ErrorPool) addError(data *errorData, args []string) {
	if data == nil {
		panic(errors.New("Error data must not be nil."))
	}

	var finalData errorData
	if args != nil && len(args) > 0 {
		finalData = p.precompileError(*data, args)
	} else {
		finalData = *data
	}

	p.errorStack.push(finalData)
}

func (p *ErrorPool) precompileError(data errorData, args []string) errorData {
	newMsg := p.replaceMasks(data.message, args)
	return errorData{name: data.name, message: newMsg, code: data.code, module: data.module, fix: data.fix}
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

func (p *ErrorPool) notify(level ErrorLevel) error {
	if !p.HasErrors() {
		panic(errors.New("Cannot notify when there are no errors."))
	}

	last := p.errorStack.peek()

	p.validateData(level, last)

	callerLine, callerFile, callerFuncName, _ := p.getErrorInfo(true)
	errorLine, errorFile, errorFuncName, errorSite := p.getErrorInfo(false)

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
			Level:          level,
			Timestamp:      time.Now(),
		}
		o.OnUpdate(info)
	}

	return nil
}

func (p *ErrorPool) validateData(level ErrorLevel, data errorData) {
	if level == Fatal || level == Error || level == Warning {
		p.validateComplete(data)
	} else {
		p.validateIncomplete(data)
	}
}

func (p *ErrorPool) validateComplete(data errorData) {
	p.checkName(data.name)
	p.checkMessage(data.message)
	p.checkCode(data.code)
	p.checkModule(data.module)
	p.checkFix(data.fix)
}

func (p *ErrorPool) validateIncomplete(data errorData) {
	p.checkMessage(data.message)
	p.checkModule(data.module)
}

func (p *ErrorPool) checkName(input string) {
	if isEmpty(input) {
		panic(errors.New("Error name must not be empty."))
	}

	if !isAlphaSpace(input) {
		panic(errors.New("Error name must only contain letters and spaces."))
	}

	if hasTrailingOrLeading(input) {
		panic(errors.New("Error name must not contain leading or trailing spaces."))
	}

	if hasConsecutiveSpaces(input) {
		panic(errors.New("Error name must not contain consecutive spaces."))
	}

	if !isSentence(input) {
		panic(errors.New(
			"Error name must follow the sentence case convention (first letter uppercase and all others lowercase).",
		))
	}
}

func (p *ErrorPool) checkMessage(input string) {
	if isEmpty(input) {
		panic(errors.New("Error message must not be empty."))
	}

	if !containsAlpha(input) {
		panic(errors.New("Error message must contain at least one letter."))
	}

	if containsControl(input) {
		panic(errors.New("Error message must not contain control characters."))
	}

	if !isCharUpper(input[0]) {
		panic(errors.New("Error message must start with uppercase letter."))
	}
}

func (p *ErrorPool) checkCode(input string) {
	if isEmpty(input) {
		panic(errors.New("Error code must not be empty."))
	}

	if len(input) != 3 {
		panic(errors.New("Error code must have 3 characters of length."))
	}

	if !isAlphaNum(input) {
		panic(errors.New("Error code must only contain uppercase letters and numbers."))
	}

	if containsLower(input) {
		panic(errors.New("Error code must contain only uppercase letters."))
	}
}

func (p *ErrorPool) checkModule(input string) {
	if isEmpty(input) {
		panic(errors.New("Error module must not be empty."))
	}

	if hasTrailingOrLeading(input) {
		panic(errors.New("Error module must not contain leading or trailing spaces."))
	}

	if hasConsecutiveSpaces(input) {
		panic(errors.New("Error module must not contain consecutive spaces."))
	}

	if !isAlphaSpace(input) {
		panic(errors.New("Error module must only contain letters and spaces."))
	}

	if !isTitle(input) {
		panic(errors.New("Error module must follow the title case convention (first letter of each word in uppercase separated by a space)."))
	}
}

func (p *ErrorPool) checkFix(input string) {
	if isEmpty(input) {
		panic(errors.New("Error fix must not be empty."))
	}

	if !containsAlpha(input) {
		panic(errors.New("Error fix must contain at least one letter."))
	}

	if containsControl(input) {
		panic(errors.New("Error fix must not contain control characters."))
	}

	if !isCharUpper(input[0]) {
		panic(errors.New("Error fix must start with uppercase letter."))
	}
}

func (p *ErrorPool) getErrorInfo(caller bool) (int, string, string, string) {
	index := p.getStackIndex(caller)

	pc, file, line, ok := runtime.Caller(index)
	funcName := p.extractFuncName(runtime.FuncForPC(pc).Name())

	if !ok {
		panic(errors.New("Could not retrieve runtime caller."))
	}

	return line, filepath.Base(file), funcName, p.getLineContents(file, line)
}

func (p *ErrorPool) getStackIndex(caller bool) int {
	if caller {
		return 2
	}
	return 3
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
		panic(errors.New("Could not open file."))
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
		panic(errors.New("Scanner could not read the file."))
	}

	panic(errors.New("File line does not exist."))
}

func (p *ErrorPool) ClearErrors() {
	p.errorStack.clear()
}

func (p *ErrorPool) HasErrors() bool {
	return !p.errorStack.isEmpty()
}

func (p *ErrorPool) GetErrorCount() int {
	return p.errorStack.size()
}

func NewErrorData(
	name string,
	message string,
	code string,
	module string,
	fix string,
) *errorData {
	return &errorData{
		name:    name,
		message: message,
		code:    code,
		module:  module,
		fix:     fix,
	}
}
