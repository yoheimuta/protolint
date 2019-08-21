package customrules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/go-protoparser/parser/meta"
	"github.com/yoheimuta/protolint/linter/report"
)

// SimpleRule verifies that all enum names are LowerSnakeCase.
type SimpleRule struct {
	verbose bool
	fixMode bool
}

// NewSimpleRule creates a new SimpleRule.
func NewSimpleRule(
	verbose bool,
	fixMode bool,
) SimpleRule {
	return SimpleRule{
		verbose: verbose,
		fixMode: fixMode,
	}
}

// ID returns the ID of this rule.
func (r SimpleRule) ID() string {
	return "SIMPLE"
}

// Purpose returns the purpose of this rule.
func (r SimpleRule) Purpose() string {
	return "Simple custom rule."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r SimpleRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r SimpleRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	return []report.Failure{
		report.Failuref(meta.Position{}, r.ID(), "Custom Rule, verbose=%v, fixMode=%v", r.verbose, r.fixMode),
	}, nil
}
