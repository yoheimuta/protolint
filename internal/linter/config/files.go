package config

import "github.com/yoheimuta/protolint/internal/stringsutil"

// Files represents the target files.
type Files struct {
	Exclude []string `yaml:"exclude"`
}

func (d Files) shouldSkipRule(
	displayPath string,
) bool {
	return stringsutil.ContainsCrossPlatformPathInSlice(displayPath, d.Exclude)
}
