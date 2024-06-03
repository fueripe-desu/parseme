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
