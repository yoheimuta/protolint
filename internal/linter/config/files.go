package config

// Files represents the target files.
type Files struct {
	Exclude []string `yaml:"exclude"`
}

func (d Files) shouldSkipRule(
	displayPath string,
) bool {
	for _, exclude := range d.Exclude {
		if displayPath == exclude {
			return true
		}
	}
	return false
}
