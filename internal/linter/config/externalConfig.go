package config

// Lint represents the lint configuration.
type Lint struct {
	Ignores     Ignores
	Files       Files
	Directories Directories
	Rules       Rules
	RulesOption RulesOption `yaml:"rules_option" json:"rules_option"`
}

type embeddedConfig struct {
	Protolint *Lint `json:"protolint"`
}

// ExternalConfig represents the external configuration.
type ExternalConfig struct {
	SourcePath string
	Lint       Lint
}

// ShouldSkipRule checks whether to skip applying the rule to the file.
func (c ExternalConfig) ShouldSkipRule(
	ruleID string,
	displayPath string,
	defaultRuleIDs []string,
) bool {
	lint := c.Lint
	return lint.Ignores.shouldSkipRule(ruleID, displayPath) ||
		lint.Files.shouldSkipRule(displayPath) ||
		lint.Directories.shouldSkipRule(displayPath) ||
		lint.Rules.shouldSkipRule(ruleID, defaultRuleIDs)
}

func (p embeddedConfig) toExternalConfig() *ExternalConfig {
	if p.Protolint == nil {
		return nil
	}

	return &ExternalConfig{
		Lint: *p.Protolint,
	}
}
