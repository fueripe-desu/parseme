package parseme

import (
	"os"
	"strings"
)

func fetchFileContents(filepath string) (*[]byte, error) {
	file, size, fileErr := openFile(filepath)

	if fileErr != nil {
		return nil, fileErr
	}

	bytes, readErr := readFile(file, size)

	if readErr != nil {
		return nil, readErr
	}

	return bytes, nil
}

func openFile(filepath string) (*os.File, int64, error) {
	fileInfo, err := os.Stat(filepath)

	if err != nil {
		return nil, 0, &FileNotFoundError{Filepath: filepath}
	}

	if fileInfo.IsDir() {
		return nil, 0, &FileIsDirError{Filepath: filepath}
	}

	file, err := os.Open(filepath)

	if err != nil && strings.Contains(err.Error(), "permission denied") {
		return nil, 0, &FileNotReadable{Filepath: filepath}
	}

	if err != nil {
		return nil, 0, err
	}

	return file, fileInfo.Size(), nil
}

func readFile(file *os.File, size int64) (*[]byte, error) {
	defer file.Close()
	_, err := file.Seek(0, 0)

	if size < 0 {
		return nil, &ReadNegativeSizeError{}
	}

	if err != nil {
		return nil, err
	}

	bytes := make([]byte, size)
	_, readErr := file.Read(bytes)

	if readErr != nil {
		return nil, readErr
	}

	return &bytes, nil
}
