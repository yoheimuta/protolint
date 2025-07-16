package util_test

import (
	"os"
	"strings"

	"github.com/yoheimuta/protolint/internal/osutil"
)

// TestData is a wrapped test file.
type TestData struct {
	FilePath   string
	OriginData []byte
}

// NewTestData create a new TestData.
func NewTestData(
	filePath string,
) (TestData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return TestData{}, nil
	}
	return TestData{
		FilePath:   filePath,
		OriginData: data,
	}, nil
}

// Data returns a content.
func (d TestData) Data() ([]byte, error) {
	return os.ReadFile(d.FilePath)
}

// Restore writes the original content back to the file.
func (d TestData) Restore() error {
	newlineChar := "\n"
	lines := strings.Split(string(d.OriginData), newlineChar)
	return osutil.WriteLinesToExistingFile(d.FilePath, lines, newlineChar)
}
