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

func TestFieldNamesLowerSnakeCaseRule_Apply(t *testing.T) {
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
			name: "no failures for proto with valid field names",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName: "song_name",
							},
							&parser.Field{
								FieldName: "singer",
							},
							&parser.Field{
								FieldName: "number2_should_be_valid",
							},
							&parser.MapField{
								MapName: "song_name2",
							},
							&parser.Oneof{
								OneofFields: []*parser.OneofField{
									{
										FieldName: "song_name3",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "failures for proto with invalid field names",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName: "song_Name",
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
								FieldName: "number_1_should_be_invalid",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   130,
										Line:     6,
										Column:   10,
									},
								},
							},
							&parser.Field{
								FieldName: "_leading_under_score_should_be_invalid",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   160,
										Line:     7,
										Column:   10,
									},
								},
							},
							&parser.Field{
								FieldName: "trailing_under_score_should_be_invalid_",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   190,
										Line:     8,
										Column:   10,
									},
								},
							},
							&parser.Field{
								FieldName: "double__under__score__should__be__invalid_",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   200,
										Line:     9,
										Column:   10,
									},
								},
							},
							&parser.MapField{
								MapName: "MapFieldName",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   210,
										Line:     14,
										Column:   30,
									},
								},
							},
							&parser.Oneof{
								OneofFields: []*parser.OneofField{
									{
										FieldName: "OneofFieldName",
										Meta: meta.Meta{
											Pos: meta.Position{
												Filename: "example.proto",
												Offset:   300,
												Line:     21,
												Column:   45,
											},
										},
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
					"FIELD_NAMES_LOWER_SNAKE_CASE",
					string(rule.SeverityError),
					`Field name "song_Name" must be underscore_separated_names like "song_name"`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   130,
						Line:     6,
						Column:   10,
					},
					"FIELD_NAMES_LOWER_SNAKE_CASE",
					string(rule.SeverityError),
					`Field name "number_1_should_be_invalid" must be underscore_separated_names like "number1_should_be_invalid"`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   160,
						Line:     7,
						Column:   10,
					},
					"FIELD_NAMES_LOWER_SNAKE_CASE",
					string(rule.SeverityError),
					`Field name "_leading_under_score_should_be_invalid" must be underscore_separated_names like "leading_under_score_should_be_invalid"`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   190,
						Line:     8,
						Column:   10,
					},
					"FIELD_NAMES_LOWER_SNAKE_CASE",
					string(rule.SeverityError),
					`Field name "trailing_under_score_should_be_invalid_" must be underscore_separated_names like "trailing_under_score_should_be_invalid"`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     9,
						Column:   10,
					},
					"FIELD_NAMES_LOWER_SNAKE_CASE",
					string(rule.SeverityError),
					`Field name "double__under__score__should__be__invalid_" must be underscore_separated_names like "double_under_score_should_be_invalid"`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   210,
						Line:     14,
						Column:   30,
					},
					"FIELD_NAMES_LOWER_SNAKE_CASE",
					string(rule.SeverityError),
					`Field name "MapFieldName" must be underscore_separated_names like "map_field_name"`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   300,
						Line:     21,
						Column:   45,
					},
					"FIELD_NAMES_LOWER_SNAKE_CASE",
					string(rule.SeverityError),
					`Field name "OneofFieldName" must be underscore_separated_names like "oneof_field_name"`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewFieldNamesLowerSnakeCaseRule(rule.SeverityError, false, autodisable.Noop)

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

func TestFieldNamesLowerSnakeCaseRule_Apply_fix(t *testing.T) {
	tests := []struct {
		name          string
		inputFilename string
		wantFilename  string
	}{
		{
			name:          "no fix for a correct proto",
			inputFilename: "lower_snake_case.proto",
			wantFilename:  "lower_snake_case.proto",
		},
		{
			name:          "fix for an incorrect proto",
			inputFilename: "invalid.proto",
			wantFilename:  "lower_snake_case.proto",
		},
		{
			name:          "fix for invalid underscore usage in proto",
			inputFilename: "invalid_underscores.proto",
			wantFilename:  "fixed_underscores.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewFieldNamesLowerSnakeCaseRule(rule.SeverityError, true, autodisable.Noop)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}

func TestFieldNamesLowerSnakeCaseRule_Apply_disable(t *testing.T) {
	tests := []struct {
		name               string
		inputFilename      string
		inputPlacementType autodisable.PlacementType
		wantFilename       string
	}{
		{
			name:          "do nothing in case of no violations",
			inputFilename: "lower_snake_case.proto",
			wantFilename:  "lower_snake_case.proto",
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
			r := rules.NewFieldNamesLowerSnakeCaseRule(rule.SeverityError, true, test.inputPlacementType)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}
