package config

// ExternalConfig represents the external configuration.
type ExternalConfig struct {
	Lint struct {
		Ignores []struct {
			ID    string   `yaml:"id"`
			Files []string `yaml:"files"`
		}
		Rules struct {
			NoDefault bool     `yaml:"no_default"`
			Add       []string `yaml:"add"`
			Remove    []string `yaml:"remove"`
		}
	}
}
