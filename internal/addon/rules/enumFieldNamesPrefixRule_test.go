package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/protolint/internal/addon/rules"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
	"github.com/yoheimuta/protolint/linter/report"
)

func TestEnumFieldNamesPrefixRule_Apply(t *testing.T) {
	tests := []struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		{
			name: "no failures for proto without enum fields",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumName: "FooBar",
					},
				},
			},
		},
		{
			name: "no failures for proto with valid enum field names",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Enum{
						EnumName: "FooBar",
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "FOO_BAR_UNSPECIFIED",
								Number: "0",
							},
							&parser.EnumField{
								Ident:  "FOO_BAR_FIRST_VALUE",
								Number: "1",
							},
							&parser.EnumField{
								Ident:  "FOO_BAR_SECOND_VALUE",
								Number: "2",
							},
						},
					},
				},
			},
		},
		{
			name: "no failures for proto with valid enum field names even when its enum name is snake case",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Enum{
						EnumName: "foo_bar",
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "FOO_BAR_UNSPECIFIED",
								Number: "0",
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
						EnumName: "FooBar",
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "BAR_UNSPECIFIED",
								Number: "0",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   100,
										Line:     5,
										Column:   10,
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
					"ENUM_FIELD_NAMES_PREFIX",
					`EnumField name "BAR_UNSPECIFIED" should have the prefix "FOO_BAR"`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewEnumFieldNamesPrefixRule(false)

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

func TestEnumFieldNamesPrefixRule_Apply_fix(t *testing.T) {
	tests := []struct {
		name          string
		inputFilename string
		wantFilename  string
	}{
		{
			name:          "no fix for a correct proto",
			inputFilename: "prefix.proto",
			wantFilename:  "prefix.proto",
		},
		{
			name:          "fix for an incorrect proto",
			inputFilename: "invalid.proto",
			wantFilename:  "prefix.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewEnumFieldNamesPrefixRule(true)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}
