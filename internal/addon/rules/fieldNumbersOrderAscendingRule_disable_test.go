package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
)

// TestFieldNumbersOrderAscendingRule_Apply_disableComments checks that the per-line
// disable comments are honoured by this rule.
// See https://github.com/yoheimuta/protolint/issues/558
func TestFieldNumbersOrderAscendingRule_Apply_disableComments(t *testing.T) {
	disableThis := func() *parser.Comment {
		return &parser.Comment{Raw: "// protolint:disable:this FIELD_NUMBERS_ORDER_ASCENDING"}
	}

	tests := []struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		{
			name: "no failures for out-of-order fields disabled by an inline comment",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName:     "a",
								FieldNumber:   "2",
								InlineComment: disableThis(),
							},
							&parser.Field{
								FieldName:     "b",
								FieldNumber:   "1",
								InlineComment: disableThis(),
							},
						},
					},
				},
			},
		},
		{
			name: "no failures for out-of-order enum fields disabled by an inline comment",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:         "A",
								Number:        "2",
								InlineComment: disableThis(),
							},
							&parser.EnumField{
								Ident:         "B",
								Number:        "1",
								InlineComment: disableThis(),
							},
						},
					},
				},
			},
		},
		{
			name: "a field following disabled ones is still checked against them",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName:     "a",
								FieldNumber:   "2",
								InlineComment: disableThis(),
							},
							&parser.Field{
								FieldName:     "b",
								FieldNumber:   "1",
								InlineComment: disableThis(),
							},
							&parser.Field{
								FieldName:   "c",
								FieldNumber: "1",
							},
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{},
					"FIELD_NUMBERS_ORDER_ASCENDING",
					string(rule.SeverityError),
					"fields %s and %s have the same number %d",
					"b", "c", 1,
				),
			},
		},
		{
			name: "a field without the comment is still checked",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName:   "a",
								FieldNumber: "2",
							},
							&parser.Field{
								FieldName:   "b",
								FieldNumber: "1",
							},
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{},
					"FIELD_NUMBERS_ORDER_ASCENDING",
					string(rule.SeverityError),
					"field %s should be after %s (ascending order expected)",
					"a", "b",
				),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewFieldNumbersOrderAscendingRule(rule.SeverityError)

			got, err := r.Apply(test.inputProto)
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}

			if !reflect.DeepEqual(got, test.wantFailures) {
				t.Errorf("got %v, but want %v", got, test.wantFailures)
			}
		})
	}
}
