package config

// RPCsHaveCommentOption represents the option for the RPCS_HAVE_COMMENT rule.
type RPCsHaveCommentOption struct {
	ShouldFollowGolangStyle bool `yaml:"should_follow_golang_style"`
}
