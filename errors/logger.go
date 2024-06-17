package errors

import "errors"

type logger struct {
	pool *ErrorPool
}

func (l *logger) Fatal(data *errorData, args []string) {
	l.pool.error(Fatal, data, args)
}

func (l *logger) Error(data *errorData, args []string) {
	l.pool.error(Error, data, args)
}

func (l *logger) Warning(data *errorData, args []string) {
	l.pool.error(Warning, data, args)
}

func (l *logger) Info(content string, module string) {
	errorData := &errorData{message: content, module: module}
	l.pool.error(Info, errorData, nil)
}

func (l *logger) Debug(content string, module string) {
	errorData := &errorData{message: content, module: module}
	l.pool.error(Debug, errorData, nil)
}

func (l *logger) Trace(content string, module string) {
	errorData := &errorData{message: content, module: module}
	l.pool.error(Trace, errorData, nil)
}

var loggerInstance *logger

func InitLogger(pool *ErrorPool) {
	if pool == nil {
		panic(errors.New("Cannot init logger if pool is nil."))
	}

	loggerInstance = &logger{pool: pool}
}

func GetLogger() *logger {
	if loggerInstance == nil {
		panic(errors.New("Must initialize logger before accessing it."))
	}

	return loggerInstance
}

func IsLoggerInitialized() bool {
	return loggerInstance != nil
}

func CloseLogger() {
	loggerInstance.pool.ClearErrors()
	loggerInstance.pool.UnsubscribeAll()
}
