package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/report"
)

func TestFieldsHaveCommentRule_Apply(t *testing.T) {
	tests := []struct {
		name                         string
		inputProto                   *parser.Proto
		inputShouldFollowGolangStyle bool
		wantFailures                 []report.Failure
	}{
		{
			name: "no failures for proto without field",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{},
			},
		},
		{
			name: "no failures for proto including valid fields with comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName: "FieldName",
								Comments: []*parser.Comment{
									{
										Raw: "// a field name.",
									},
								},
							},
							&parser.MapField{
								MapName: "MapFieldName",
								Comments: []*parser.Comment{
									{
										Raw: "// a map field name.",
									},
								},
							},
							&parser.Oneof{
								OneofFields: []*parser.OneofField{
									{
										FieldName: "OneofFieldName",
										Comments: []*parser.Comment{
											{
												Raw: "// a oneof field name.",
											},
										},
									},
								},
							},
							&parser.Field{
								FieldName: "FieldName",
								InlineComment: &parser.Comment{
									Raw: "// a field name.",
								},
							},
							&parser.MapField{
								MapName: "MapFieldName",
								InlineComment: &parser.Comment{
									Raw: "// a map field name.",
								},
							},
							&parser.Oneof{
								OneofFields: []*parser.OneofField{
									{
										FieldName: "OneofFieldName",
										InlineComment: &parser.Comment{
											Raw: "// a oneof field name.",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "no failures for proto including valid fields with Golang style comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName: "FieldName",
								Comments: []*parser.Comment{
									{
										Raw: "// FieldName is a field name.",
									},
								},
							},
							&parser.MapField{
								MapName: "MapFieldName",
								Comments: []*parser.Comment{
									{
										Raw: "// MapFieldName is a map field name.",
									},
								},
							},
							&parser.Oneof{
								OneofFields: []*parser.OneofField{
									{
										FieldName: "OneofFieldName",
										Comments: []*parser.Comment{
											{
												Raw: "// OneofFieldName is a oneof field name.",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			inputShouldFollowGolangStyle: true,
		},
		{
			name: "failures for proto with invalid fields",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName: "FieldName",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   150,
										Line:     7,
										Column:   15,
									},
								},
							},
							&parser.MapField{
								MapName: "MapFieldName",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   200,
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
						Offset:   150,
						Line:     7,
						Column:   15,
					},
					"FIELDS_HAVE_COMMENT",
					`Field "FieldName" should have a comment`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     14,
						Column:   30,
					},
					"FIELDS_HAVE_COMMENT",
					`Field "MapFieldName" should have a comment`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   300,
						Line:     21,
						Column:   45,
					},
					"FIELDS_HAVE_COMMENT",
					`Field "OneofFieldName" should have a comment`,
				),
			},
		},
		{
			name: "failures for proto with invalid fields without Golang style comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName: "FieldName",
								Comments: []*parser.Comment{
									{
										Raw: "// a field name.",
									},
								},
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   150,
										Line:     7,
										Column:   15,
									},
								},
							},
							&parser.MapField{
								MapName: "MapFieldName",
								Comments: []*parser.Comment{
									{
										Raw: "// a map field name.",
									},
								},
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   200,
										Line:     14,
										Column:   30,
									},
								},
							},
							&parser.Oneof{
								OneofFields: []*parser.OneofField{
									{
										FieldName: "OneofFieldName",
										Comments: []*parser.Comment{
											{
												Raw: "// a oneof field name.",
											},
										},
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
			inputShouldFollowGolangStyle: true,
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   150,
						Line:     7,
						Column:   15,
					},
					"FIELDS_HAVE_COMMENT",
					`Field "FieldName" should have a comment of the form "// FieldName ..."`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     14,
						Column:   30,
					},
					"FIELDS_HAVE_COMMENT",
					`Field "MapFieldName" should have a comment of the form "// MapFieldName ..."`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   300,
						Line:     21,
						Column:   45,
					},
					"FIELDS_HAVE_COMMENT",
					`Field "OneofFieldName" should have a comment of the form "// OneofFieldName ..."`,
				),
			},
		},
		{
			name: "failures for proto with fields without Golang style comments due to the inline comment",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Field{
								FieldName: "FieldName2",
								InlineComment: &parser.Comment{
									Raw: "// FieldName2 is a field name.",
								},
							},
							&parser.MapField{
								MapName: "MapFieldName2",
								InlineComment: &parser.Comment{
									Raw: "// MapFieldName2 is a map field name.",
								},
							},
							&parser.Oneof{
								OneofFields: []*parser.OneofField{
									{
										FieldName: "OneofFieldName2",
										InlineComment: &parser.Comment{
											Raw: "// OneofFieldName2 is a oneof field name.",
										},
									},
								},
							},
						},
					},
				},
			},
			inputShouldFollowGolangStyle: true,
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{},
					"FIELDS_HAVE_COMMENT",
					`Field "FieldName2" should have a comment of the form "// FieldName2 ..."`,
				),
				report.Failuref(
					meta.Position{},
					"FIELDS_HAVE_COMMENT",
					`Field "MapFieldName2" should have a comment of the form "// MapFieldName2 ..."`,
				),
				report.Failuref(
					meta.Position{},
					"FIELDS_HAVE_COMMENT",
					`Field "OneofFieldName2" should have a comment of the form "// OneofFieldName2 ..."`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewFieldsHaveCommentRule(test.inputShouldFollowGolangStyle)

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
