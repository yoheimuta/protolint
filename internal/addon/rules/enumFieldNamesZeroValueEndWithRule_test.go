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

func TestEnumFieldNamesZeroValueEndWithRule_Apply(t *testing.T) {
	tests := []struct {
		name         string
		inputProto   *parser.Proto
		inputSuffix  string
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
								Ident:  "ENUM_TYPE_UNSPECIFIED",
								Number: "0",
							},
							&parser.EnumField{
								Ident:  "SECOND_VALUE",
								Number: "1",
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
								Ident:  "FIRST_VALUE",
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
					"ENUM_FIELD_NAMES_ZERO_VALUE_END_WITH",
					string(rule.SeverityError),
					`EnumField name "FIRST_VALUE" with zero value should have the suffix "UNSPECIFIED"`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewEnumFieldNamesZeroValueEndWithRule(rule.SeverityError, test.inputSuffix, false, autodisable.Noop)

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

func TestEnumFieldNamesZeroValueEndWithRule_Apply_fix(t *testing.T) {
	tests := []struct {
		name          string
		inputFilename string
		wantFilename  string
	}{
		{
			name:          "no fix for a correct proto",
			inputFilename: "suffix.proto",
			wantFilename:  "suffix.proto",
		},
		{
			name:          "fix for an incorrect proto",
			inputFilename: "invalid.proto",
			wantFilename:  "suffix.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewEnumFieldNamesZeroValueEndWithRule(rule.SeverityError, "", true, autodisable.Noop)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}

func TestEnumFieldNamesZeroValueEndWithRule_Apply_disable(t *testing.T) {
	tests := []struct {
		name               string
		inputFilename      string
		inputPlacementType autodisable.PlacementType
		wantFilename       string
	}{
		{
			name:          "do nothing in case of no violations",
			inputFilename: "suffix.proto",
			wantFilename:  "suffix.proto",
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
			r := rules.NewEnumFieldNamesZeroValueEndWithRule(rule.SeverityError, "", true, test.inputPlacementType)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}
