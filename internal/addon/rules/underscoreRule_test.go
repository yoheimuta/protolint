package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
)

func TestUnderscoreRule_Apply(t *testing.T) {
	tests := []struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		{
			name: "no failures for proto without fields",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{},
				},
			},
		},
		{
			name: "no failures for the names the guide allows",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{FieldName: "number2"},
							&parser.Field{FieldName: "number_two"},
							&parser.MapField{MapName: "labels_by_key"},
							&parser.Oneof{
								OneofFields: []*parser.OneofField{
									{FieldName: "choice_name"},
								},
							},
							&parser.Enum{
								EnumBody: []parser.Visitee{
									&parser.EnumField{Ident: "UNSPECIFIED"},
									&parser.EnumField{Ident: "VALUE_A"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "failures for every shape the guide forbids",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName: "number_1_should_be_invalid",
								Meta:      meta.Meta{Pos: meta.Position{Offset: 10, Line: 2, Column: 1}},
							},
							&parser.Field{
								FieldName: "_leading_should_be_invalid",
								Meta:      meta.Meta{Pos: meta.Position{Offset: 20, Line: 3, Column: 1}},
							},
							&parser.Field{
								FieldName: "trailing_should_be_invalid_",
								Meta:      meta.Meta{Pos: meta.Position{Offset: 30, Line: 4, Column: 1}},
							},
							&parser.Field{
								FieldName: "double__under__score",
								Meta:      meta.Meta{Pos: meta.Position{Offset: 40, Line: 5, Column: 1}},
							},
							&parser.Field{
								FieldName: "address_line_1",
								Meta:      meta.Meta{Pos: meta.Position{Offset: 45, Line: 5, Column: 20}},
							},
							&parser.MapField{
								MapName: "map_name_1",
								Meta:    meta.Meta{Pos: meta.Position{Offset: 50, Line: 6, Column: 1}},
							},
							&parser.Oneof{
								OneofFields: []*parser.OneofField{
									{
										FieldName: "oneof_name_2",
										Meta:      meta.Meta{Pos: meta.Position{Offset: 60, Line: 7, Column: 1}},
									},
								},
							},
							&parser.Enum{
								EnumBody: []parser.Visitee{
									&parser.EnumField{
										Ident: "VALUE_1",
										Meta:  meta.Meta{Pos: meta.Position{Offset: 70, Line: 8, Column: 1}},
									},
								},
							},
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(meta.Position{Offset: 10, Line: 2, Column: 1}, "UNDERSCORE", string(rule.SeverityWarning), `Field name "number_1_should_be_invalid" must have every underscore followed by a letter`),
				report.Failuref(meta.Position{Offset: 20, Line: 3, Column: 1}, "UNDERSCORE", string(rule.SeverityWarning), `Field name "_leading_should_be_invalid" must not begin with an underscore`),
				report.Failuref(meta.Position{Offset: 30, Line: 4, Column: 1}, "UNDERSCORE", string(rule.SeverityWarning), `Field name "trailing_should_be_invalid_" must not end with an underscore`),
				report.Failuref(meta.Position{Offset: 40, Line: 5, Column: 1}, "UNDERSCORE", string(rule.SeverityWarning), `Field name "double__under__score" must have every underscore followed by a letter`),
				report.Failuref(meta.Position{Offset: 45, Line: 5, Column: 20}, "UNDERSCORE", string(rule.SeverityWarning), `Field name "address_line_1" must have every underscore followed by a letter`),
				report.Failuref(meta.Position{Offset: 50, Line: 6, Column: 1}, "UNDERSCORE", string(rule.SeverityWarning), `Map field name "map_name_1" must have every underscore followed by a letter`),
				report.Failuref(meta.Position{Offset: 60, Line: 7, Column: 1}, "UNDERSCORE", string(rule.SeverityWarning), `Oneof field name "oneof_name_2" must have every underscore followed by a letter`),
				report.Failuref(meta.Position{Offset: 70, Line: 8, Column: 1}, "UNDERSCORE", string(rule.SeverityWarning), `Enum field name "VALUE_1" must have every underscore followed by a letter`),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, err := rules.NewUnderscoreRule(rule.SeverityWarning, autodisable.Noop).Apply(test.inputProto)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if !reflect.DeepEqual(got, test.wantFailures) {
				t.Errorf("got = %v, but want = %v", got, test.wantFailures)
			}
		})
	}
}
