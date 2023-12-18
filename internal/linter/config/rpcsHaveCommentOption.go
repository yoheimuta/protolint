package config

// RPCsHaveCommentOption represents the option for the RPCS_HAVE_COMMENT rule.
type RPCsHaveCommentOption struct {
	CustomizableSeverityOption
	ShouldFollowGolangStyle bool `yaml:"should_follow_golang_style" json:"should_follow_golang_style" toml:"should_follow_golang_style"`
}
