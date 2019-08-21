package rule

import "github.com/yoheimuta/protolint/linter/rule"

// Rules is a list of Rules.
type Rules []rule.Rule

// Default returns a default set of rules.
func (rs Rules) Default() Rules {
	var d Rules
	for _, r := range rs {
		if r.IsOfficial() {
			d = append(d, r)
		}
	}
	return d
}

// IDs returns a set of rule ids.
func (rs Rules) IDs() []string {
	return ruleIDs(rs)
}

func ruleIDs(rules []rule.Rule) []string {
	var ids []string
	for _, rule := range rules {
		ids = append(ids, rule.ID())
	}
	return ids
}
