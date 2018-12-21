package rules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolinter/internal/addon/rules/internal/visitor"
	"github.com/yoheimuta/protolinter/internal/linter/report"
	"github.com/yoheimuta/protolinter/internal/strs"
)

// MessageNamesUpperCamelCaseRule verifies that all message names are CamelCase (with an initial capital).
// See https://developers.google.com/protocol-buffers/docs/style#message-and-field-names.
type MessageNamesUpperCamelCaseRule struct{}

// NewMessageNamesUpperCamelCaseRule creates a new MessageNamesUpperCamelCaseRule.
func NewMessageNamesUpperCamelCaseRule() MessageNamesUpperCamelCaseRule {
	return MessageNamesUpperCamelCaseRule{}
}

// ID returns the ID of this rule.
func (r MessageNamesUpperCamelCaseRule) ID() string {
	return "MESSAGE_NAMES_UPPER_CAMEL_CASE"
}

// Purpose returns the purpose of this rule.
func (r MessageNamesUpperCamelCaseRule) Purpose() string {
	return "Verifies that all message names are CamelCase (with an initial capital)."
}

// Apply applies the rule to the proto.
func (r MessageNamesUpperCamelCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &messageNamesUpperCamelCaseVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type messageNamesUpperCamelCaseVisitor struct {
	*visitor.BaseAddVisitor
}

// VisitMessage checks the message.
func (v *messageNamesUpperCamelCaseVisitor) VisitMessage(message *parser.Message) bool {
	if !strs.IsUpperCamelCase(message.MessageName) {
		v.AddFailuref(message.Meta.Pos, "Message name %q must be UpperCamelCase", message.MessageName)
	}
	return true
}
