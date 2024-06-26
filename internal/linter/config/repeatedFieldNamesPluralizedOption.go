package config

// RepeatedFieldNamesPluralizedOption represents the option for the REPEATED_FIELD_NAMES_PLURALIZED rule.
type RepeatedFieldNamesPluralizedOption struct {
	CustomizableSeverityOption `yaml:",inline"`
	PluralRules                map[string]string `yaml:"plural_rules" json:"plural_rules" toml:"plural_rules"`
	SingularRules              map[string]string `yaml:"singular_rules" json:"singular_rules" toml:"singular_rules"`
	UncountableRules           []string          `yaml:"uncountable_rules" json:"uncountable_rules" toml:"uncountable_rules"`
	IrregularRules             map[string]string `yaml:"irregular_rules" json:"irregular_rules" toml:"irregular_rules"`
}
