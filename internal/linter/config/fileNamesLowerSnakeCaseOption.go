package config

// FileNamesLowerSnakeCaseOption represents the option for the FILE_NAMES_LOWER_SNAKE_CASE rule.
type FileNamesLowerSnakeCaseOption struct {
	CustomizableSeverityOption `yaml:",inline"`
	Excludes                   []string `yaml:"excludes" json:"excludes" toml:"excludes"`
}
