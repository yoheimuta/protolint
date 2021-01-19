package config

import (
	"github.com/yoheimuta/protolint/internal/filepathutil"
)

// Directories represents the target directories.
type Directories struct {
	Exclude []string `yaml:"exclude"`
}

func (d Directories) shouldSkipRule(
	displayPath string,
) bool {
	for _, exclude := range d.Exclude {
		if filepathutil.HasUnixPathPrefix(displayPath,
			exclude+string(filepathutil.OSPathSeparator)) {
			return true
		}
	}
	return false
}
