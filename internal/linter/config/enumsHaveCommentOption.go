package config

// EnumsHaveCommentOption represents the option for the ENUMS_HAVE_COMMENT rule.
type EnumsHaveCommentOption struct {
	CustomizableSeverityOption `yaml:",inline"`
	ShouldFollowGolangStyle    bool `yaml:"should_follow_golang_style" json:"should_follow_golang_style" toml:"should_follow_golang_style"`
}
