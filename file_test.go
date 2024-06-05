package parseme

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_openFile(t *testing.T) {
	openFunc := func(filepath string) (*os.File, int64) {
		fileInfo, _ := os.Stat(filepath)
		size := fileInfo.Size()
		file, _ := os.Open(filepath)

		return file, size
	}

	// Tests for the file return value from openFile()
	testcasesFile := []struct {
		name        string
		filepath    string
		expected    func(filepath string) (*os.File, int64)
		expectedErr error
	}{
		{
			"open file",
			"test_data/example1.html",
			openFunc,
			nil,
		},
		{
			"open unicode file",
			"test_data/example2.html",
			openFunc,
			nil,
		},
		{
			"unknown file path",
			"unknown/example.html",
			openFunc,
			&FileNotFoundError{Filepath: "unknown/example.html"},
		},
		{
			"dir path",
			"./",
			openFunc,
			&FileIsDirError{Filepath: "./"},
		},
	}

	// Tests for the size return value from openFile()
	testcasesSize := []struct {
		name        string
		filepath    string
		expected    int64
		expectedErr error
	}{
		{
			"regular size file",
			"test_data/example1.html",
			141,
			nil,
		},
		{
			"empty file",
			"test_data/example3.html",
			0,
			nil,
		},
	}

	for _, tc := range testcasesFile {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			file, size, err := openFile(tc.filepath)

			if tc.expectedErr != nil {
				assert.Equal(tc.expectedErr, err)
				assert.EqualError(err, tc.expectedErr.Error())
			} else if err != nil && tc.expectedErr == nil {
				t.Log("openFile() returned an error, but expected error is null")
				t.Log(err.Error())
				t.FailNow()
			} else {
				expectedFile, expectedSize := tc.expected(tc.filepath)
				assert.Equal((*file).Name(), (*expectedFile).Name())
				assert.Equal(size, expectedSize)
			}
		})
	}

	for _, tc := range testcasesSize {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			_, size, err := openFile(tc.filepath)

			if tc.expectedErr != nil {
				assert.Equal(tc.expectedErr, err)
				assert.EqualError(err, tc.expectedErr.Error())
			} else if err != nil && tc.expectedErr == nil {
				t.Log("openFile() returned an error, but expected error is null")
				t.Log(err.Error())
				t.FailNow()
			} else {
				assert.Equal(size, tc.expected)
			}
		})
	}

	t.Run("file without read permission", func(t *testing.T) {
		assert := assert.New(t)
		filepath := "test_data/temp_permission.html"
		file, createErr := os.Create(filepath)

		t.Cleanup(func() {
			file.Close()
			os.Chmod(filepath, 0777)
			os.Remove(filepath)
		})

		if createErr != nil {
			t.Log("Could not create temporary file:", createErr.Error())
			t.FailNow()
		}

		bytes := []byte{60, 104, 116, 109, 108, 62}
		_, writeErr := file.Write(bytes)

		if writeErr != nil {
			t.Log("Could not write to temporary file:", writeErr.Error())
			t.FailNow()
		}

		os.Chmod(filepath, 0111)

		_, _, openErr := openFile(filepath)

		if openErr == nil {
			t.Log("Should not be able to read file without permission.")
			t.FailNow()
		} else {
			expectedErr := &FileNotReadable{Filepath: filepath}
			assert.Equal(openErr, expectedErr)
			assert.EqualError(openErr, expectedErr.Error())
		}
	})
}

func Test_readFile(t *testing.T) {
	testcasesFile := []struct {
		name     string
		filepath string
		expected func() *[]byte
	}{
		{
			"read file",
			"test_data/example1.html",
			func() *[]byte {
				bytes, _ := os.ReadFile("test_data/example1.html")
				return &bytes
			},
		},
		{
			"read file with unicode chars",
			"test_data/example2.html",
			func() *[]byte {
				bytes, _ := os.ReadFile("test_data/example2.html")
				return &bytes
			},
		},
	}

	testcasesSize := []struct {
		name        string
		filepath    string
		size        int64
		expected    []byte
		expectedErr error
	}{
		{
			"negative size as parameter",
			"test_data/example1.html",
			-1,
			[]byte{},
			&ReadNegativeSizeError{},
		},
		{
			"zero size as parameter",
			"test_data/example1.html",
			0,
			[]byte{},
			nil,
		},
		{
			"size less than file length",
			"test_data/example1.html",
			2,
			[]byte{60, 104},
			nil,
		},
	}

	for _, tc := range testcasesFile {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			expected := tc.expected()
			file, size, _ := openFile(tc.filepath)
			result, _ := readFile(file, size)
			assert.Equal(expected, result)
		})
	}

	for _, tc := range testcasesSize {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			file, _, _ := openFile(tc.filepath)
			result, err := readFile(file, tc.size)

			if tc.expectedErr != nil {
				assert.Equal(err, tc.expectedErr)
				assert.EqualError(err, tc.expectedErr.Error())
			} else if err != nil && tc.expectedErr == nil {
				t.Log("readFile() returned an error, but expected error is null")
				t.FailNow()
			} else {
				assert.Equal(tc.expected, *result)
			}
		})
	}

}
