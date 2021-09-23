package osutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Constants used by LineEnding.
const (
	lf   = "\n"
	cr   = "\r"
	crlf = "\r\n"
)

// ReadAllLines reads all lines from a file.
func ReadAllLines(
	fileName string,
	newlineChar string,
) ([]string, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), newlineChar), nil
}

// WriteLinesToExistingFile writes lines to an existing file.
func WriteLinesToExistingFile(
	fileName string,
	lines []string,
	newlineChar string,
) error {
	data := strings.Join(lines, newlineChar)
	return WriteExistingFile(
		fileName,
		[]byte(data),
	)
}

// WriteExistingFile writes the byte array to an existing file.
func WriteExistingFile(
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

// DetectLineEnding detects a dominant line ending in the content.
func DetectLineEnding(content string) (string, error) {
	prev := ' '
	counts := make(map[string]int)
	for _, c := range content {
		if c == '\r' && prev != '\n' {
			counts[cr]++
		} else if c == '\n' {
			if prev == '\r' {
				counts[crlf]++
				counts[cr]--
			} else {
				counts[lf]++
			}
		}
		prev = c
	}
	if counts[crlf]+counts[cr]+counts[lf] == 0 {
		return "", nil
	}

	if counts[crlf] > counts[cr] && counts[crlf] > counts[lf] {
		return crlf, nil
	} else if counts[cr] > counts[lf] && counts[cr] > counts[crlf] {
		return cr, nil
	} else if counts[lf] > counts[cr] && counts[lf] > counts[crlf] {
		return lf, nil
	}
	return "", fmt.Errorf("not found dominant line ending, counts=%v", counts)
}
