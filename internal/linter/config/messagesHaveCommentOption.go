package config

// MessagesHaveCommentOption represents the option for the MESSAGES_HAVE_COMMENT rule.
type MessagesHaveCommentOption struct {
	CustomizableSeverityOption
	ShouldFollowGolangStyle bool `yaml:"should_follow_golang_style"`
}
