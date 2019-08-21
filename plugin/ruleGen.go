package plugin

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/internal/linter/report"
	"github.com/yoheimuta/protolint/internal/linter/rule"
)

// RuleGen is a generator for a rule. It's adapted to rule.Rule interface.
type RuleGen func(
	verbose bool,
	fixMode bool,
) rule.Rule

// ID implements rule.Rule.
func (RuleGen) ID() string {
	return ""
}

// Purpose implements rule.Rule.
func (RuleGen) Purpose() string {
	return ""
}

// IsOfficial implements rule.Rule.
func (RuleGen) IsOfficial() bool {
	return true
}

// Apply implements rule.Rule.
func (RuleGen) Apply(proto *parser.Proto) ([]report.Failure, error) {
	return nil, nil
}
