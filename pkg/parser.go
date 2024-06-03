package parseme

import (
	"os"
)

type HtmlParser struct {
	filepath string
}

func NewHtmlParser(filepath string) *HtmlParser {
	return &HtmlParser{filepath: filepath}
}

func (p *HtmlParser) Parse() (*[]byte, error) {
	file, size, fileErr := openFile(p.filepath)

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

	if err != nil {
		return nil, 0, &FileError{Message: err.Error() + ":", Filepath: filepath}
	}

	return file, fileInfo.Size(), nil
}

func readFile(file *os.File, size int64) (*[]byte, error) {
	defer file.Close()
	_, err := file.Seek(0, 0)

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
