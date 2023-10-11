package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/protolint/internal/addon/rules"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
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
			name: "no failures for proto with valid enum field names even when its camel-case string starts with a 2-letter abbreviation",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Enum{
						EnumName: "ITDepartmentRegion",
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "IT_DEPARTMENT_REGION_UNSPECIFIED",
								Number: "0",
							},
						},
					},
				},
			},
		},
		{
			name: "no failures for proto with valid enum field names even when its camel-case string includes OAuth",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Enum{
						EnumName: "ListAccountPipedriveOAuthsEnabledFilter",
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident:  "LIST_ACCOUNT_PIPEDRIVE_OAUTHS_ENABLED_FILTER_UNSPECIFIED",
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
			rule := rules.NewEnumFieldNamesPrefixRule(rule.SeverityError, false, autodisable.Noop)

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
			r := rules.NewEnumFieldNamesPrefixRule(rule.SeverityError, true, autodisable.Noop)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}

func TestEnumFieldNamesPrefixRule_Apply_disable(t *testing.T) {
	tests := []struct {
		name               string
		inputFilename      string
		inputPlacementType autodisable.PlacementType
		wantFilename       string
	}{
		{
			name:          "do nothing in case of no violations",
			inputFilename: "prefix.proto",
			wantFilename:  "prefix.proto",
		},
		{
			name:               "insert disable:next comments",
			inputFilename:      "invalid.proto",
			inputPlacementType: autodisable.Next,
			wantFilename:       "disable_next.proto",
		},
		{
			name:               "insert disable:this comments",
			inputFilename:      "invalid.proto",
			inputPlacementType: autodisable.ThisThenNext,
			wantFilename:       "disable_this.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewEnumFieldNamesPrefixRule(rule.SeverityError, true, test.inputPlacementType)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}
