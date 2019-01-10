package config

// Rules represents the enabled rule set.
type Rules struct {
	NoDefault bool     `yaml:"no_default"`
	Add       []string `yaml:"add"`
	Remove    []string `yaml:"remove"`
}
