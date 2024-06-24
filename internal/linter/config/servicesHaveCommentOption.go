package config

// ServicesHaveCommentOption represents the option for the SERVICES_HAVE_COMMENT rule.
type ServicesHaveCommentOption struct {
	CustomizableSeverityOption `yaml:",inline"`
	ShouldFollowGolangStyle    bool `yaml:"should_follow_golang_style" json:"should_follow_golang_style" toml:"should_follow_golang_style"`
}
