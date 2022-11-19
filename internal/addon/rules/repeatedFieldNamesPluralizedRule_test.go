package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/report"
)

func TestRepeatedFieldNamesPluralizedRule_Apply(t *testing.T) {
	tests := []struct {
		name             string
		pluralRules      map[string]string
		singularRules    map[string]string
		uncountableRules []string
		irregularRules   map[string]string
		inputProto       *parser.Proto
		wantFailures     []report.Failure
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
			name: "no failures for proto with valid field names",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName: "singer",
							},
							&parser.Field{
								IsRepeated: true,
								FieldName:  "singers",
							},
							&parser.GroupField{
								IsRepeated: true,
								GroupName:  "people",
								MessageBody: []parser.Visitee{
									&parser.Field{
										IsRepeated: true,
										FieldName:  "some_singers",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "no failures for proto with valid field names considering the rule is case insensitive",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								IsRepeated: true,
								FieldName:  "RunningOnDeviceIDS",
							},
							&parser.GroupField{
								IsRepeated: true,
								GroupName:  "RunningOnDeviceIDs",
							},
						},
					},
				},
			},
		},
		{
			name: "no failures for proto with field names by applying some customization",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								IsRepeated: true,
								FieldName:  "regexii",
							},
							&parser.Field{
								IsRepeated: true,
								FieldName:  "paper",
							},
							&parser.Field{
								IsRepeated: true,
								FieldName:  "paper",
							},
							&parser.Field{
								IsRepeated: true,
								FieldName:  "regular",
							},
						},
					},
				},
			},
			pluralRules: map[string]string{
				"(?i)gex$": "gexii",
			},
			singularRules: map[string]string{
				"(?i)gexii": "gex",
			},
			uncountableRules: []string{
				"paper",
			},
			irregularRules: map[string]string{
				"irregular": "regular",
			},
		},
		{
			name: "failures for proto with non-pluralized repeated field names",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								IsRepeated: true,
								FieldName:  "singer",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   100,
										Line:     5,
										Column:   10,
									},
								},
							},
							&parser.Field{
								IsRepeated: true,
								FieldName:  "persons",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   200,
										Line:     10,
										Column:   20,
									},
								},
							},
							&parser.GroupField{
								IsRepeated: true,
								GroupName:  "some_singer",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   210,
										Line:     14,
										Column:   30,
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
					"REPEATED_FIELD_NAMES_PLURALIZED",
					`Repeated field name "singer" must be pluralized name "singers"`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"REPEATED_FIELD_NAMES_PLURALIZED",
					`Repeated field name "persons" must be pluralized name "people"`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   210,
						Line:     14,
						Column:   30,
					},
					"REPEATED_FIELD_NAMES_PLURALIZED",
					`Repeated group name "some_singer" must be pluralized name "some_singers"`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewRepeatedFieldNamesPluralizedRule(
				test.pluralRules,
				test.singularRules,
				test.uncountableRules,
				test.irregularRules,
				false,
				autodisable.Noop,
			)

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

func TestRepeatedFieldNamesPluralizedRule_Apply_fix(t *testing.T) {
	tests := []struct {
		name             string
		pluralRules      map[string]string
		singularRules    map[string]string
		uncountableRules []string
		irregularRules   map[string]string
		inputFilename    string
		wantFilename     string
	}{
		{
			name:          "no fix for a correct proto",
			inputFilename: "pluralized.proto",
			wantFilename:  "pluralized.proto",
		},
		{
			name:          "fix for an incorrect proto",
			inputFilename: "invalid.proto",
			wantFilename:  "pluralized.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewRepeatedFieldNamesPluralizedRule(
				test.pluralRules,
				test.singularRules,
				test.uncountableRules,
				test.irregularRules,
				true,
				autodisable.Noop,
			)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}

func TestRepeatedFieldNamesPluralizedRule_Apply_disable(t *testing.T) {
	tests := []struct {
		name               string
		pluralRules        map[string]string
		singularRules      map[string]string
		uncountableRules   []string
		irregularRules     map[string]string
		inputFilename      string
		inputPlacementType autodisable.PlacementType
		wantFilename       string
	}{
		{
			name:          "do nothing in case of no violations",
			inputFilename: "pluralized.proto",
			wantFilename:  "pluralized.proto",
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
			r := rules.NewRepeatedFieldNamesPluralizedRule(
				test.pluralRules,
				test.singularRules,
				test.uncountableRules,
				test.irregularRules,
				true,
				test.inputPlacementType,
			)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}
