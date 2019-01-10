package config

// Ignore represents files ignoring the specific rule.
type Ignore struct {
	ID    string   `yaml:"id"`
	Files []string `yaml:"files"`
}
