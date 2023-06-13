package plugin

import (
	"github.com/yoheimuta/protolint/internal/addon/plugin/proto"
	"github.com/yoheimuta/protolint/internal/addon/plugin/shared"

	"github.com/yoheimuta/protolint/linter/rule"
)

// GetExternalRules provides the external rules.
func GetExternalRules(
	clients []shared.RuleSet,
	fixMode bool,
	verbose bool,
) ([]rule.Rule, error) {
	var rs []rule.Rule

	for _, client := range clients {
		resp, err := client.ListRules(&proto.ListRulesRequest{
			Verbose: verbose,
			FixMode: fixMode,
		})
		if err != nil {
			return nil, err
		}

		for _, r := range resp.Rules {
			severity := getSeverity(r.Severity)
			rs = append(rs, newExternalRule(r.Id, r.Purpose, client, severity))
		}
	}
	return rs, nil
}

func getSeverity(ruleSeverity proto.RuleSeverity) rule.Severity {
	switch ruleSeverity {
	case proto.RuleSeverity_RULE_SEVERITY_ERROR:
		return rule.SeverityError
	case proto.RuleSeverity_RULE_SEVERITY_WARNING:
		return rule.SeverityWarning
	case proto.RuleSeverity_RULE_SEVERITY_NOTE:
		return rule.SeverityNote
	}

	return rule.SeverityError
}
