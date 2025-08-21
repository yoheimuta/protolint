package rules

import (
	"strconv"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// FieldNumbersOrderAscendingRule verifies the order of fields.
type FieldNumbersOrderAscendingRule struct {
	RuleWithSeverity
}

// NewFieldNumbersOrderAscendingRule creates a new FieldsOrderAscendingRule.
func NewFieldNumbersOrderAscendingRule(severity rule.Severity) FieldNumbersOrderAscendingRule {
	return FieldNumbersOrderAscendingRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
	}
}

// ID returns the ID of this rule.
func (r FieldNumbersOrderAscendingRule) ID() string {
	return "FIELD_NUMBERS_ORDER_ASCENDING"
}

// Purpose returns the purpose of this rule.
func (r FieldNumbersOrderAscendingRule) Purpose() string {
	return "Verifies the order of fields."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r FieldNumbersOrderAscendingRule) IsOfficial() bool {
	return false
}

// Apply applies the rule to the proto.
func (r FieldNumbersOrderAscendingRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &fieldNumbersOrderAscendingVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID(), string(r.Severity())),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type fieldNumbersOrderAscendingVisitor struct {
	*visitor.BaseAddVisitor
}

// VisitMessage checks the message
func (v *fieldNumbersOrderAscendingVisitor) VisitMessage(message *parser.Message) bool {
	var (
		lastNumber int
		lastName   string
		hasError   bool
	)

	for _, element := range message.MessageBody {
		field, ok := element.(*parser.Field)
		if !ok {
			continue
		}

		number, err := strconv.Atoi(field.FieldNumber)
		if err != nil {
			v.AddFailuref(
				field.Meta.Pos,
				"field number '%s' is not a number",
				field.FieldNumber,
			)

			hasError = true
			continue
		}

		if number <= 0 {
			v.AddFailuref(
				field.Meta.Pos,
				"field number should be positive",
			)

			hasError = true
			continue
		}

		number, isError := v.isAscending(field.Meta.Pos, field.FieldName, number, lastName, lastNumber)
		if isError {
			hasError = true
		}

		lastNumber = number
		lastName = field.FieldName
	}

	return !hasError
}

// VisitEnum checks the enum
func (v *fieldNumbersOrderAscendingVisitor) VisitEnum(enum *parser.Enum) bool {
	var (
		lastNumber int = -1
		lastIdent  string
		hasError   bool
	)

	for _, element := range enum.EnumBody {
		field, ok := element.(*parser.EnumField)
		if !ok {
			continue
		}

		number, err := strconv.Atoi(field.Number)
		if err != nil {
			v.AddFailuref(
				field.Meta.Pos,
				"field number '%s' is not a number",
				field.Number,
			)

			hasError = true
			continue
		}

		if number < 0 {
			v.AddFailuref(
				field.Meta.Pos,
				"field number should be positive",
			)

			hasError = true
			continue
		}

		number, isError := v.isAscending(field.Meta.Pos, field.Ident, number, lastIdent, lastNumber)
		if isError {
			hasError = true
		}

		lastNumber = number
		lastIdent = field.Ident
	}

	return !hasError
}

func (v *fieldNumbersOrderAscendingVisitor) isAscending(
	pos meta.Position, fieldName string, number int, lastName string, lastNumber int,
) (curNumber int, hasError bool) {
	if number == lastNumber {
		v.AddFailuref(
			pos,
			"fields %s and %s have the same number %d",
			lastName, fieldName, number,
		)

		return number, true
	}

	if number < lastNumber {
		v.AddFailuref(
			pos,
			"field %s should be after %s (ascending order expected)",
			lastName, fieldName,
		)

		return number, true
	}

	return number, false
}
