package setting_test

import (
	"os"
	"path/filepath"
)

func projectRootPath() string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "yoheimuta", "protolint")
}

// TestDataPath is the directory path for the data used by tests.
func TestDataPath(elem ...string) string {
	ps := []string{
		projectRootPath(),
		"_testdata",
	}
	ps = append(ps, elem...)
	return filepath.Join(ps...)
}
