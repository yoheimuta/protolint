package config

// MaxLineLengthOption represents the option for the MAX_LINE_LENGTH rule.
type MaxLineLengthOption struct {
	CustomizableSeverityOption
	MaxChars int `yaml:"max_chars"`
	TabChars int `yaml:"tab_chars"`
}
