package errors

import "errors"

type logger struct {
	pool   *ErrorPool
	module string
}

func (l *logger) Fatal(data *errorData, args []string) {
	l.log(Fatal, l.module, data, args)
}

func (l *logger) Error(data *errorData, args []string) {
	l.log(Error, l.module, data, args)
}

func (l *logger) Warning(data *errorData, args []string) {
	l.log(Warning, l.module, data, args)
}

func (l *logger) Info(content string, module string) {
	l.log(Info, l.module, &errorData{message: content}, nil)
}

func (l *logger) Debug(content string, module string) {
	l.log(Debug, l.module, &errorData{message: content}, nil)
}

func (l *logger) Trace(content string, module string) {
	l.log(Trace, l.module, &errorData{message: content}, nil)
}

func (l *logger) log(
	level ErrorLevel,
	module string,
	data *errorData,
	args []string,
) {
	if l.pool != nil {
		l.pool.error(level, module, data, args)
	}
}

var loggerInstance *logger

func InitLogger(pool *ErrorPool) {
	loggerInstance = &logger{pool: pool}
}

func GetLogger() *logger {
	if !IsLoggerInitialized() {
		panic(errors.New("Must initialize logger before accessing it."))
	}

	return loggerInstance
}

func IsLoggerInitialized() bool {
	return loggerInstance != nil
}

func CloseLogger() {
	if IsLoggerInitialized() && loggerInstance.pool != nil {
		loggerInstance.pool.ClearErrors()
		loggerInstance.pool.UnsubscribeAll()
	}
}

func SetLoggerModule(module string) {
	if !IsLoggerInitialized() {
		panic(errors.New("Cannot set module if logger is not initialized."))
	}

	loggerInstance.module = module
}
