package config

import (
	"fmt"
	"strings"
)

// MaxLineLengthOption represents the option for the MAX_LINE_LENGTH rule.
type MaxLineLengthOption struct {
	MaxChars int `yaml:"max_chars"`
	TabChars int `yaml:"tab_chars"`
}

// IndentOption represents the option for the INDENT rule.
type IndentOption struct {
	Style   string
	Newline string
}

// UnmarshalYAML implements yaml.v2 Unmarshaler interface.
func (i *IndentOption) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var option struct {
		Style   string `yaml:"style"`
		Newline string `yaml:"newline"`
	}
	if err := unmarshal(&option); err != nil {
		return err
	}

	var style string
	switch option.Style {
	case "tab":
		style = "\t"
	case "4":
		style = strings.Repeat(" ", 4)
	case "2":
		style = strings.Repeat(" ", 2)
	case "":
		break
	default:
		return fmt.Errorf("%s is an invalid style option. valid option is tab, 4 or 2", option.Style)
	}
	i.Style = style

	switch option.Newline {
	case "\n", "\r", "\r\n", "":
		i.Newline = option.Newline
	default:
		return fmt.Errorf("%s is an invalid newline option. valid option is \n, \r or \r\n", option.Newline)
	}
	return nil
}

// RulesOption represents the option for some rules.
type RulesOption struct {
	MaxLineLength MaxLineLengthOption `yaml:"max_line_length"`
	Indent        IndentOption        `yaml:"indent"`
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
		RulesOption RulesOption `yaml:"rules_option"`
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
