package config

// ExternalConfig represents the external configuration.
type ExternalConfig struct {
	Lint struct {
		Ignores     Ignores
		Rules       Rules
		RulesOption RulesOption `yaml:"rules_option"`
	}
}

// ShouldSkipRule checks whether to skip applying the rule to the file.
func (c ExternalConfig) ShouldSkipRule(
	ruleID string,
	displayPath string,
	defaultRuleIDs []string,
) bool {
	lint := c.Lint
	return lint.Ignores.shouldSkipRule(ruleID, displayPath) ||
		lint.Rules.shouldSkipRule(ruleID, defaultRuleIDs)
}
