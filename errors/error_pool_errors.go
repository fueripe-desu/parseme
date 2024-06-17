package errors

type ObserverNotFoundError struct{}

func (e *ObserverNotFoundError) Error() string {
	return "Observer does not exist."
}

type ObserverDuplicateError struct{}

func (e *ObserverDuplicateError) Error() string {
	return "Observer already exists."
}

type NotifyObserverError struct{}

func (e *NotifyObserverError) Error() string {
	return "Cannot notify when there are no errors."
}

var nilErrorDataError *errorData = NewErrorData(
	"Nil error data",
	"Error data must not be nil.",
	"ED1",
	"Error Pool",
	"Try handling nil error data pointers before passing it to the function, or assing a new address to error data.",
)
