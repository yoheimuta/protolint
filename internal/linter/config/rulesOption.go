package config

// RulesOption represents the option for some rules.
type RulesOption struct {
	MaxLineLength MaxLineLengthOption `yaml:"max_line_length"`
	Indent        IndentOption        `yaml:"indent"`
}
