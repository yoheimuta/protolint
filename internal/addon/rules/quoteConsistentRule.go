package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/internal/linter/config"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// QuoteConsistentRule verifies that the use of quote for strings is consistent.
type QuoteConsistentRule struct {
	RuleWithSeverity
	quote config.QuoteType

	fixMode bool
}

// NewQuoteConsistentRule creates a new QuoteConsistentRule.
func NewQuoteConsistentRule(
	severity rule.Severity,
	quote config.QuoteType,
	fixMode bool,
) QuoteConsistentRule {
	return QuoteConsistentRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		quote:            quote,
		fixMode:          fixMode,
	}
}

// ID returns the ID of this rule.
func (r QuoteConsistentRule) ID() string {
	return "QUOTE_CONSISTENT"
}

// Purpose returns the purpose of this rule.
func (r QuoteConsistentRule) Purpose() string {
	return "Verifies that the use of quote for strings is consistent."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r QuoteConsistentRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r QuoteConsistentRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), r.fixMode, proto, string(r.Severity()))
	if err != nil {
		return nil, err
	}

	v := &quoteConsistentVisitor{
		BaseFixableVisitor: base,
		quote:              r.quote,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type quoteConsistentVisitor struct {
	*visitor.BaseFixableVisitor
	quote config.QuoteType
}

func (v quoteConsistentVisitor) VisitSyntax(s *parser.Syntax) bool {
	str := s.ProtobufVersionQuote
	converted := convertConsistentQuote(str, v.quote)
	if str != converted {
		v.AddFailuref(s.Meta.Pos, "Quoted string should be %s but was %s.", converted, str)
		v.Fixer.ReplaceText(s.Meta.Pos.Line, str, converted)
	}
	return false
}

func (v quoteConsistentVisitor) VisitImport(i *parser.Import) (next bool) {
	str := i.Location
	converted := convertConsistentQuote(str, v.quote)
	if str != converted {
		v.AddFailuref(i.Meta.Pos, "Quoted string should be %s but was %s.", converted, str)
		v.Fixer.ReplaceText(i.Meta.Pos.Line, str, converted)
	}
	return false
}

func (v quoteConsistentVisitor) VisitOption(o *parser.Option) (next bool) {
	str := o.Constant
	converted := convertConsistentQuote(str, v.quote)
	if str != converted {
		v.AddFailuref(o.Meta.Pos, "Quoted string should be %s but was %s.", converted, str)
		v.Fixer.ReplaceText(o.Meta.Pos.Line, str, converted)
	}
	return false
}

func (v quoteConsistentVisitor) VisitEnumField(f *parser.EnumField) (next bool) {
	for _, option := range f.EnumValueOptions {
		str := option.Constant
		converted := convertConsistentQuote(str, v.quote)
		if str != converted {
			v.AddFailuref(f.Meta.Pos, "Quoted string should be %s but was %s.", converted, str)
			v.Fixer.ReplaceText(f.Meta.Pos.Line, str, converted)
		}
	}
	return false
}

func (v quoteConsistentVisitor) VisitField(f *parser.Field) (next bool) {
	for _, option := range f.FieldOptions {
		str := option.Constant
		converted := convertConsistentQuote(str, v.quote)
		if str != converted {
			v.AddFailuref(f.Meta.Pos, "Quoted string should be %s but was %s.", converted, str)
			v.Fixer.ReplaceText(f.Meta.Pos.Line, str, converted)
		}
	}
	return false
}

func convertConsistentQuote(s string, quote config.QuoteType) string {
	var valid, invalid string
	if quote == config.DoubleQuote {
		valid, invalid = "\"", "'"
	} else {
		valid, invalid = "'", "\""
	}

	if strings.HasPrefix(s, invalid) && strings.HasSuffix(s, invalid) {
		return valid + s[1:len(s)-1] + valid
	}
	return s
}
