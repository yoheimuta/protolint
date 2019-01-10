package osutil

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// ReadAllLines reads all lines from a file.
func ReadAllLines(
	fileName string,
	newLineChar string,
) ([]string, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), newLineChar), nil
}

// WriteLinesToExistingFile writes lines to an existing file.
func WriteLinesToExistingFile(
	fileName string,
	lines []string,
	newLineChar string,
) error {
	data := strings.Join(lines, newLineChar)
	return writeExistingFile(
		fileName,
		[]byte(data),
	)
}

func writeExistingFile(
	fileName string,
	data []byte,
) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
