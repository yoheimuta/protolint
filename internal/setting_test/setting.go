package setting_test

import (
	"path/filepath"
	"runtime"
)

func projectRootPath() string {
	_, this, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Dir(filepath.Dir(filepath.Dir(this)))
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
