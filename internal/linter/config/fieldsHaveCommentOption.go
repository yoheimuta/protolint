package config

// FieldsHaveCommentOption represents the option for the FIELDS_HAVE_COMMENT rule.
type FieldsHaveCommentOption struct {
	CustomizableSeverityOption `yaml:",inline"`
	ShouldFollowGolangStyle    bool `yaml:"should_follow_golang_style" json:"should_follow_golang_style" toml:"should_follow_golang_style"`
}
