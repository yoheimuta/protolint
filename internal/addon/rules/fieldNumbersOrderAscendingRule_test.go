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

func TestFieldNumbersOrderAscendingRule_Apply(t *testing.T) {
	tests := []struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		{
			name: "no failures for proto enum with allow_alias and adjacent duplicate numbers",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.Option{OptionName: "allow_alias", Constant: "true"},
							&parser.EnumField{Ident: "FIRST_VALUE", Number: "1"},
							&parser.EnumField{Ident: "FIRST_VALUE_ALIAS", Number: "1"},
						},
					},
				},
			},
			wantFailures: nil,
		},
		{
			name: "failures for proto enum with allow_alias but non-adjacent decrease (1,2,1)",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.Option{OptionName: "allow_alias", Constant: "true"},
							&parser.EnumField{Ident: "FIRST_VALUE", Number: "1"},
							&parser.EnumField{Ident: "SECOND_VALUE", Number: "2"},
							&parser.EnumField{Ident: "THIRD_VALUE", Number: "1"},
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{},
					"FIELD_NUMBERS_ORDER_ASCENDING",
					string(rule.SeverityError),
					"field SECOND_VALUE should be after THIRD_VALUE (ascending order expected)",
				),
			},
		},
		{
			name: "no failures for proto enum with allow_alias and zero value alias (0,0)",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.Option{OptionName: "allow_alias", Constant: "true"},
							&parser.EnumField{Ident: "VALUE_UNSPECIFIED", Number: "0"},
							&parser.EnumField{Ident: "VALUE_UNSPECIFIED_ALIAS", Number: "0"},
						},
					},
				},
			},
			wantFailures: nil,
		},
		{
			name: "no failures for proto enum with ascending order numbers",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "FIRST_VALUE",
								Number: "1",
							},
							&parser.EnumField{
								Ident:  "SECOND_VALUE",
								Number: "2",
							},
						},
					},
				},
			},
			wantFailures: nil,
		},
		{
			name: "no failures for proto enum started from 0",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "VALUE_UNSPECIFIED",
								Number: "0",
							},
							&parser.EnumField{
								Ident:  "FIRST_VALUE",
								Number: "1",
							},
						},
					},
				},
			},
			wantFailures: nil,
		},
		{
			name: "no failures for proto enum with ascending order numbers with gap",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "FIRST_VALUE",
								Number: "1",
							},
							&parser.EnumField{
								Ident:  "SECOND_VALUE",
								Number: "3",
							},
						},
					},
				},
			},
			wantFailures: nil,
		},
		{
			name: "failures for proto enum with negative number",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "FIRST_VALUE",
								Number: "-1",
							},
							&parser.EnumField{
								Ident:  "SECOND_VALUE",
								Number: "2",
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
					"field number should be positive integer",
				),
			},
		},
		{
			name: "failures for proto enum with duplicated numbers",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "FIRST_VALUE",
								Number: "1",
							},
							&parser.EnumField{
								Ident:  "SECOND_VALUE",
								Number: "1",
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
					"fields FIRST_VALUE and SECOND_VALUE have the same number 1",
				),
			},
		},
		{
			name: "failures for proto enum with descending numbers",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "SECOND_VALUE",
								Number: "2",
							},
							&parser.EnumField{
								Ident:  "FIRST_VALUE",
								Number: "1",
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
					"field SECOND_VALUE should be after FIRST_VALUE (ascending order expected)",
				),
			},
		},
		{
			name: "no failures for proto message with ascending order numbers",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName:   "FIRST_VALUE",
								FieldNumber: "1",
							},
							&parser.Field{
								FieldName:   "SECOND_VALUE",
								FieldNumber: "2",
							},
						},
					},
				},
			},
			wantFailures: nil,
		},
		{
			name: "no failures for proto message with ascending order numbers with gap",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName:   "FIRST_VALUE",
								FieldNumber: "1",
							},
							&parser.Field{
								FieldName:   "SECOND_VALUE",
								FieldNumber: "3",
							},
						},
					},
				},
			},
			wantFailures: nil,
		},
		{
			name: "failures for proto message with duplicated numbers",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName:   "FIRST_VALUE",
								FieldNumber: "1",
							},
							&parser.Field{
								FieldName:   "SECOND_VALUE",
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
					"fields FIRST_VALUE and SECOND_VALUE have the same number 1",
				),
			},
		},
		{
			name: "failures for proto message with descending numbers",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName:   "SECOND_VALUE",
								FieldNumber: "2",
							},
							&parser.Field{
								FieldName:   "FIRST_VALUE",
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
					"field SECOND_VALUE should be after FIRST_VALUE (ascending order expected)",
				),
			},
		},
		{
			name: "failures for proto message with not number value",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName:   "FIRST_VALUE",
								FieldNumber: "one",
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
					"field number 'one' is not a number",
				),
			},
		},
		{
			name: "failures for proto message with not ascending multiply values",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName:   "FIRST_VALUE",
								FieldNumber: "1",
							},
							&parser.Field{
								FieldName:   "FOURTH_VALUE",
								FieldNumber: "4",
							},
							&parser.Field{
								FieldName:   "THIRD_VALUE",
								FieldNumber: "3",
							},
							&parser.Field{
								FieldName:   "SECOND_VALUE",
								FieldNumber: "2",
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
					"field FOURTH_VALUE should be after THIRD_VALUE (ascending order expected)",
				),
				report.Failuref(
					meta.Position{},
					"FIELD_NUMBERS_ORDER_ASCENDING",
					string(rule.SeverityError),
					"field THIRD_VALUE should be after SECOND_VALUE (ascending order expected)",
				),
			},
		},
		{
			name: "no failures for proto message with enum inside",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Enum{
								EnumName: "MY_ENUM",
								EnumBody: []parser.Visitee{
									&parser.Service{},
									&parser.EnumField{
										Ident:  "FIRST_ENUM_VALUE",
										Number: "0",
									},
								},
								Comments:                     nil,
								InlineComment:                nil,
								InlineCommentBehindLeftCurly: nil,
								Meta:                         meta.Meta{},
							},
							&parser.Field{
								FieldName:   "FIRST_VALUE",
								FieldNumber: "1",
							},
							&parser.Field{
								FieldName:   "SECOND_VALUE",
								FieldNumber: "2",
							},
						},
					},
				},
			},
			wantFailures: nil,
		},
		{
			name: "failures for proto message with number 0",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName:   "FIRST_VALUE",
								FieldNumber: "0",
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
					"field number should be positive integer",
				),
			},
		},
		{
			name: "failures for proto message with negative number",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName:   "FIRST_VALUE",
								FieldNumber: "-1",
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
					"field number should be positive integer",
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
