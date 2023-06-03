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
			rs = append(rs, newExternalRule(r.Id, r.Purpose, client, rule.Severity_Error))
		}
	}
	return rs, nil
}
