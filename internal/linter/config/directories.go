package config

import (
	"path/filepath"
)

// Directories represents the target directories.
type Directories struct {
	Exclude []string `yaml:"exclude"`
}

func (d Directories) shouldSkipRule(
	displayPath string,
) bool {
	displayDir := filepath.Dir(displayPath)

	for _, exclude := range d.Exclude {
		if displayDir == exclude {
			return true
		}
	}
	return false
}
