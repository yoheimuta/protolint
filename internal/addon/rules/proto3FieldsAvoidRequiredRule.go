package rules

import (
	"github.com/yoheimuta/go-protoparser/v4/lexer"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/fixer"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// Proto3FieldsAvoidRequiredRule verifies that all fields should avoid required for proto3.
// See https://developers.google.com/protocol-buffers/docs/style#things-to-avoid
type Proto3FieldsAvoidRequiredRule struct {
	RuleWithSeverity
	fixMode bool
}

// NewProto3FieldsAvoidRequiredRule creates a new Proto3FieldsAvoidRequiredRule.
func NewProto3FieldsAvoidRequiredRule(
	severity rule.Severity,
	fixMode bool,
) Proto3FieldsAvoidRequiredRule {
	return Proto3FieldsAvoidRequiredRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		fixMode:          fixMode,
	}
}

// ID returns the ID of this rule.
func (r Proto3FieldsAvoidRequiredRule) ID() string {
	return "PROTO3_FIELDS_AVOID_REQUIRED"
}

// Purpose returns the purpose of this rule.
func (r Proto3FieldsAvoidRequiredRule) Purpose() string {
	return "Verifies that all fields should avoid required for proto3."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r Proto3FieldsAvoidRequiredRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r Proto3FieldsAvoidRequiredRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	base, err := visitor.NewBaseFixableVisitor(r.ID(), r.fixMode, proto)
	if err != nil {
		return nil, err
	}

	v := &proto3FieldsAvoidRequiredVisitor{
		BaseFixableVisitor: base,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type proto3FieldsAvoidRequiredVisitor struct {
	*visitor.BaseFixableVisitor
	isProto3 bool
}

// VisitSyntax checks the syntax.
func (v *proto3FieldsAvoidRequiredVisitor) VisitSyntax(s *parser.Syntax) bool {
	v.isProto3 = s.ProtobufVersion == "proto3"
	return false
}

// VisitField checks the field.
func (v *proto3FieldsAvoidRequiredVisitor) VisitField(field *parser.Field) bool {
	if v.isProto3 && field.IsRequired {
		v.AddFailuref(field.Meta.Pos, `Field %q should avoid required for proto3`, field.FieldName)

		err := v.Fixer.SearchAndReplace(field.Meta.Pos, func(lex *lexer.Lexer) fixer.TextEdit {
			lex.NextKeyword()
			return fixer.TextEdit{
				Pos:     lex.Pos.Offset,
				End:     lex.Pos.Offset + len(lex.Text),
				NewText: []byte(""),
			}
		})
		if err != nil {
			panic(err)
		}
	}
	return false
}
