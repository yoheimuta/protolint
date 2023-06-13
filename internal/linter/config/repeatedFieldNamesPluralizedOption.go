package config

// RepeatedFieldNamesPluralizedOption represents the option for the REPEATED_FIELD_NAMES_PLURALIZED rule.
type RepeatedFieldNamesPluralizedOption struct {
	CustomizableSeverityOption
	PluralRules      map[string]string `yaml:"plural_rules"`
	SingularRules    map[string]string `yaml:"singular_rules"`
	UncountableRules []string          `yaml:"uncountable_rules"`
	IrregularRules   map[string]string `yaml:"irregular_rules"`
}
