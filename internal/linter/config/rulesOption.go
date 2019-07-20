package config

// RulesOption represents the option for some rules.
type RulesOption struct {
	MaxLineLength                         MaxLineLengthOption                   `yaml:"max_line_length"`
	Indent                                IndentOption                          `yaml:"indent"`
	ServiceNamesEndWith                   ServiceNamesEndWithOption             `yaml:"service_names_end_with"`
	FieldNamesExcludePrepositionsOption   FieldNamesExcludePrepositionsOption   `yaml:"field_names_exclude_prepositions"`
	MessageNamesExcludePrepositionsOption MessageNamesExcludePrepositionsOption `yaml:"message_names_exclude_prepositions"`
}
