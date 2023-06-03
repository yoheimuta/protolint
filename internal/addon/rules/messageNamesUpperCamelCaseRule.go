package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/fixer"
	"github.com/yoheimuta/protolint/linter/rule"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/strs"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// MessageNamesUpperCamelCaseRule verifies that all message names are CamelCase (with an initial capital).
// See https://developers.google.com/protocol-buffers/docs/style#message-and-field-names.
type MessageNamesUpperCamelCaseRule struct {
	RuleWithSeverity
	fixMode         bool
	autoDisableType autodisable.PlacementType
}

// NewMessageNamesUpperCamelCaseRule creates a new MessageNamesUpperCamelCaseRule.
func NewMessageNamesUpperCamelCaseRule(
	severity rule.Severity,
	fixMode bool,
	autoDisableType autodisable.PlacementType,
) MessageNamesUpperCamelCaseRule {
	if autoDisableType != autodisable.Noop {
		fixMode = false
	}
	return MessageNamesUpperCamelCaseRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		fixMode:          fixMode,
		autoDisableType:  autoDisableType,
	}
}

// ID returns the ID of this rule.
func (r MessageNamesUpperCamelCaseRule) ID() string {
	return "MESSAGE_NAMES_UPPER_CAMEL_CASE"
}

// Purpose returns the purpose of this rule.
func (r MessageNamesUpperCamelCaseRule) Purpose() string {
	return "Verifies that all message names are CamelCase (with an initial capital)."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r MessageNamesUpperCamelCaseRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r MessageNamesUpperCamelCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), r.fixMode, proto)
	if err != nil {
		return nil, err
	}

	v := &messageNamesUpperCamelCaseVisitor{
		BaseFixableVisitor: base,
	}
	return visitor.RunVisitorAutoDisable(v, proto, r.ID(), r.autoDisableType)
}

type messageNamesUpperCamelCaseVisitor struct {
	*visitor.BaseFixableVisitor
}

// VisitMessage checks the message.
func (v *messageNamesUpperCamelCaseVisitor) VisitMessage(message *parser.Message) bool {
	name := message.MessageName
	if !strs.IsUpperCamelCase(name) {
		expected := strs.ToUpperCamelCase(name)
		v.AddFailuref(message.Meta.Pos, "Message name %q must be UpperCamelCase like %q", name, expected)

		err := v.Fixer.SearchAndReplace(message.Meta.Pos, func(lex *lexer.Lexer) fixer.TextEdit {
			lex.NextKeyword()
			lex.Next()
			return fixer.TextEdit{
				Pos:     lex.Pos.Offset,
				End:     lex.Pos.Offset + len(lex.Text) - 1,
				NewText: []byte(expected),
			}
		})
		if err != nil {
			panic(err)
		}
	}
	return true
}
