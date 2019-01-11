package config

import (
	"os"
	"strings"
)

// Directories represents the target directories.
type Directories struct {
	Exclude []string `yaml:"exclude"`
}

func (d Directories) shouldSkipRule(
	displayPath string,
) bool {
	for _, exclude := range d.Exclude {
		if strings.HasPrefix(displayPath, exclude+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}
