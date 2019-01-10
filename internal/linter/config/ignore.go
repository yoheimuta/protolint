package config

import "github.com/yoheimuta/protolint/internal/stringsutil"

// Ignore represents files ignoring the specific rule.
type Ignore struct {
	ID    string   `yaml:"id"`
	Files []string `yaml:"files"`
}

func (i Ignore) shouldSkipRule(
	ruleID string,
	displayPath string,
) bool {
	if i.ID != ruleID {
		return false
	}
	return stringsutil.ContainsStringInSlice(displayPath, i.Files)
}
