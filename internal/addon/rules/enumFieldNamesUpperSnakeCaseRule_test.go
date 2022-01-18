package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/report"
)

func TestEnumFieldNamesUpperSnakeCaseRule_Apply(t *testing.T) {
	tests := []struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		{
			name: "no failures for proto without enum fields",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{},
				},
			},
		},
		{
			name: "no failures for proto with valid enum field names",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident: "FIRST_VALUE",
							},
							&parser.EnumField{
								Ident: "SECOND_VALUE",
							},
						},
					},
				},
			},
		},
		{
			name: "failures for proto with invalid enum field names",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident: "fIRST_VALUE",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   100,
										Line:     5,
										Column:   10,
									},
								},
							},
							&parser.EnumField{
								Ident: "secondValue",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   200,
										Line:     10,
										Column:   20,
									},
								},
							},
						},
					},
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
					`EnumField name "fIRST_VALUE" must be CAPITALS_WITH_UNDERSCORES like "FIRST_VALUE"`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
					`EnumField name "secondValue" must be CAPITALS_WITH_UNDERSCORES like "SECOND_VALUE"`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewEnumFieldNamesUpperSnakeCaseRule(false)

			got, err := rule.Apply(test.inputProto)
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

func TestEnumFieldNamesUpperSnakeCaseRule_Apply_fix(t *testing.T) {
	tests := []struct {
		name          string
		inputFilename string
		wantFilename  string
	}{
		{
			name:          "no fix for a correct proto",
			inputFilename: "upperSnakeCase.proto",
			wantFilename:  "upperSnakeCase.proto",
		},
		{
			name:          "fix for an incorrect proto",
			inputFilename: "invalid.proto",
			wantFilename:  "upperSnakeCase.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewEnumFieldNamesUpperSnakeCaseRule(true)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}
