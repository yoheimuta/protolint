package config

import (
	"strings"

	"github.com/yoheimuta/protolint/internal/filepathutil"
)

// Directories represents the target directories.
type Directories struct {
	Exclude []string `yaml:"exclude" json:"exclude" toml:"exclude"`
}

func (d Directories) shouldSkipRule(
	displayPath string,
) bool {
	for _, exclude := range d.Exclude {
		if !strings.HasSuffix(exclude, string(filepathutil.OSPathSeparator)) {
			exclude += string(filepathutil.OSPathSeparator)
		}
		if filepathutil.HasUnixPathPrefix(displayPath, exclude) {
			return true
		}
	}
	return false
}
