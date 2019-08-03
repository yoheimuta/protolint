package config

// FieldsHaveCommentOption represents the option for the FIELDS_HAVE_COMMENT rule.
type FieldsHaveCommentOption struct {
	ShouldFollowGolangStyle bool `yaml:"should_follow_golang_style"`
}
