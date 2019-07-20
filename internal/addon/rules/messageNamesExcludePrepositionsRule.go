package rules

import (
	"strings"

	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/internal/addon/rules/internal/visitor"
	"github.com/yoheimuta/protolint/internal/linter/report"
	"github.com/yoheimuta/protolint/internal/stringsutil"
	"github.com/yoheimuta/protolint/internal/strs"
)

// MessageNamesExcludePrepositionsRule verifies that all message names don't include prepositions (e.g. "With", "For").
// It is assumed that the message names are CamelCase (with an initial capital).
// See https://cloud.google.com/apis/design/naming_convention#message_names.
type MessageNamesExcludePrepositionsRule struct {
	prepositions []string
}

// NewMessageNamesExcludePrepositionsRule creates a new MessageNamesExcludePrepositionsRule.
func NewMessageNamesExcludePrepositionsRule(
	prepositions []string,
) MessageNamesExcludePrepositionsRule {
	if len(prepositions) == 0 {
		for _, p := range defaultPrepositions {
			prepositions = append(prepositions, strings.Title(p))
		}
	}
	return MessageNamesExcludePrepositionsRule{
		prepositions: prepositions,
	}
}

// ID returns the ID of this rule.
func (r MessageNamesExcludePrepositionsRule) ID() string {
	return "MESSAGE_NAMES_EXCLUDE_PREPOSITIONS"
}

// Purpose returns the purpose of this rule.
func (r MessageNamesExcludePrepositionsRule) Purpose() string {
	return `Verifies that all message names don't include prepositions (e.g. "With", "For").`
}

// Apply applies the rule to the proto.
func (r MessageNamesExcludePrepositionsRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &messageNamesExcludePrepositionsVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(),
		prepositions:   r.prepositions,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type messageNamesExcludePrepositionsVisitor struct {
	*visitor.BaseAddVisitor
	prepositions []string
}

// VisitMessage checks the message.
func (v *messageNamesExcludePrepositionsVisitor) VisitMessage(message *parser.Message) bool {
	parts := strs.SplitCamelCaseWord(message.MessageName)
	for _, p := range parts {
		if stringsutil.ContainsStringInSlice(p, v.prepositions) {
			v.AddFailuref(message.Meta.Pos, "Message name %q should not include a preposition %q", message.MessageName, p)
		}
	}
	return true
}
