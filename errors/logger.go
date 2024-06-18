package errors

import "errors"

type logger struct {
	pool   *ErrorPool
	module string
}

func (l *logger) Fatal(data *errorData, args []string) {
	l.pool.error(Fatal, l.module, data, args)
}

func (l *logger) Error(data *errorData, args []string) {
	l.pool.error(Error, l.module, data, args)
}

func (l *logger) Warning(data *errorData, args []string) {
	l.pool.error(Warning, l.module, data, args)
}

func (l *logger) Info(content string, module string) {
	l.pool.error(Info, l.module, &errorData{message: content}, nil)
}

func (l *logger) Debug(content string, module string) {
	l.pool.error(Debug, l.module, &errorData{message: content}, nil)
}

func (l *logger) Trace(content string, module string) {
	l.pool.error(Trace, l.module, &errorData{message: content}, nil)
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
	if loggerInstance != nil {
		loggerInstance.pool.ClearErrors()
		loggerInstance.pool.UnsubscribeAll()
	}
}

func SetLoggerModule(module string) {
	if loggerInstance == nil {
		panic(errors.New("Cannot set module if logger is not initialized."))
	}

	loggerInstance.module = module
}
