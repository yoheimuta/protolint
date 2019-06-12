package config

// RulesOption represents the option for some rules.
type RulesOption struct {
	MaxLineLength       MaxLineLengthOption       `yaml:"max_line_length"`
	Indent              IndentOption              `yaml:"indent"`
	ServiceNamesEndWith ServiceNamesEndWithOption `yaml:"service_names_end_with"`
}
