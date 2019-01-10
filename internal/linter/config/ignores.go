package config

// Ignores represents list about files ignoring the specific rule.
type Ignores []Ignore

func (is Ignores) shouldSkipRule(
	ruleID string,
	displayPath string,
) bool {
	for _, ignore := range is {
		if ignore.shouldSkipRule(ruleID, displayPath) {
			return true
		}
	}
	return false
}
