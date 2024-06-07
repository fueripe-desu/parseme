package parseme

import "fmt"

type FileError struct {
	Message  string
	Filepath string
}

func (e *FileError) Error() string {
	return fmt.Sprintln(e.Message, e.Filepath)
}

type FileNotFoundError struct {
	Filepath string
}

func (e *FileNotFoundError) Error() string {
	fileError := &FileError{Message: "File not found in the given path:", Filepath: e.Filepath}
	return fileError.Error()
}

type FileIsDirError struct {
	Filepath string
}

func (e *FileIsDirError) Error() string {
	fileError := &FileError{Message: "Given filepath points to a dir instead of a file:", Filepath: e.Filepath}
	return fileError.Error()
}

type FileNotReadable struct {
	Filepath string
}

func (e *FileNotReadable) Error() string {
	fileError := &FileError{Message: "File does not have read permission:", Filepath: e.Filepath}
	return fileError.Error()
}

type ReadNegativeSizeError struct{}

func (e *ReadNegativeSizeError) Error() string {
	return "Cannot read file if scale is negative."
}

type StackError struct {
	message string
}

func (e *StackError) Error() string {
	return e.message
}

// Property error
type PropertyError struct {
	Message string
}

func (e *PropertyError) Error() string {
	return e.Message
}

type PropertyBooleanValueError struct{}

func (e *PropertyBooleanValueError) Error() string {
	err := PropertyError{Message: "Cannot get boolean value of non-boolean elements."}
	return err.Error()
}

type PropertyInvalidNameError struct {
	name string
}

func (e *PropertyInvalidNameError) Error() string {
	msg := fmt.Sprintf("The property name '%v' is not valid.", e.name)
	err := PropertyError{Message: msg}
	return err.Error()
}

type PropertyEmptyNameError struct{}

func (e *PropertyEmptyNameError) Error() string {
	err := PropertyError{Message: "Name must not be empty."}
	return err.Error()
}

type PropertyInvalidBooleanError struct {
	value string
}

func (e *PropertyInvalidBooleanError) Error() string {
	msg := fmt.Sprintf("Boolean property value must be either 'true' or 'false'. Instead got: '%v'", e.value)
	err := PropertyError{Message: msg}
	return err.Error()
}
