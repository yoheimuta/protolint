package config

// MaxLineLengthOption represents the option for the MAX_LINE_LENGTH rule.
type MaxLineLengthOption struct {
	MaxChars int `yaml:"max_chars"`
	TabChars int `yaml:"tab_chars"`
}

// RuleOption represents the option for some rules.
type RuleOption struct {
	MaxLineLength MaxLineLengthOption `yaml:"max_line_length"`
}

// ExternalConfig represents the external configuration.
type ExternalConfig struct {
	Lint struct {
		Ignores []struct {
			ID    string   `yaml:"id"`
			Files []string `yaml:"files"`
		}
		Rules struct {
			NoDefault bool     `yaml:"no_default"`
			Add       []string `yaml:"add"`
			Remove    []string `yaml:"remove"`
		}
		RuleOption RuleOption `yaml:"rule_option"`
	}
}

// SkipRule checks to skip applying the rule to the file.
func (c ExternalConfig) SkipRule(
	ruleID string,
	displayPath string,
	defaultRuleIDs []string,
) bool {
	if ignoreFiles, ok := c.ignores(ruleID); ok {
		if containsStringInSlice(displayPath, ignoreFiles) {
			return true
		}
	}

	rules := c.rules(defaultRuleIDs)
	return !containsStringInSlice(ruleID, rules)
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
		if !containsStringInSlice(id, c.Lint.Rules.Remove) {
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

func containsStringInSlice(
	needle string,
	haystack []string,
) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}
