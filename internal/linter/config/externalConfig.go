package config

import "github.com/yoheimuta/protolint/internal/stringsutil"

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
	if ignoreFiles, ok := c.ignores(ruleID); ok {
		if stringsutil.ContainsStringInSlice(displayPath, ignoreFiles) {
			return true
		}
	}

	rules := c.rules(defaultRuleIDs)
	return !stringsutil.ContainsStringInSlice(ruleID, rules)
}

func (c ExternalConfig) rules(
	defaultRuleIDs []string,
) []string {
	ruleIDs := defaultRuleIDs
	if c.Lint.Rules.NoDefault {
		ruleIDs = ruleIDs[:0]
	}
	for _, add := range c.Lint.Rules.Add {
		ruleIDs = append(ruleIDs, add)
	}

	var newRuleIDs []string
	for _, id := range ruleIDs {
		if !stringsutil.ContainsStringInSlice(id, c.Lint.Rules.Remove) {
			newRuleIDs = append(newRuleIDs, id)
		}
	}
	return newRuleIDs
}

func (c ExternalConfig) ignores(
	ruleID string,
) ([]string, bool) {
	for _, ignore := range c.Lint.Ignores {
		if ignore.ID == ruleID {
			return ignore.Files, true
		}
	}
	return nil, false
}
