package config

// EnumsHaveCommentOption represents the option for the ENUMS_HAVE_COMMENT rule.
type EnumsHaveCommentOption struct {
	ShouldFollowGolangStyle bool `yaml:"should_follow_golang_style"`
}
