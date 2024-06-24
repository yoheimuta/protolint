package config

// EnumFieldNamesZeroValueEndWithOption represents the option for the ENUM_FIELD_NAMES_ZERO_VALUE_END_WITH rule.
type EnumFieldNamesZeroValueEndWithOption struct {
	CustomizableSeverityOption `yaml:",inline"`
	Suffix                     string `yaml:"suffix" json:"suffix" toml:"suffix"`
}
