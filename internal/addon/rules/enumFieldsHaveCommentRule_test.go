package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
)

func TestEnumFieldsHaveCommentRule_Apply(t *testing.T) {
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
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident: "EnumFieldName",
								Comments: []*parser.Comment{
									{
										Raw: "// a field name.",
									},
								},
							},
							&parser.EnumField{
								Ident: "EnumFieldName2",
								InlineComment: &parser.Comment{
									Raw: "// a field name.",
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
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident: "EnumFieldName",
								Comments: []*parser.Comment{
									{
										Raw: "// EnumFieldName is a field name.",
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
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident: "EnumFieldName",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   150,
										Line:     7,
										Column:   15,
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
					"ENUM_FIELDS_HAVE_COMMENT",
					`EnumField "EnumFieldName" should have a comment`,
				),
			},
		},
		{
			name: "failures for proto with invalid fields without Golang style comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident: "EnumFieldName",
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
					"ENUM_FIELDS_HAVE_COMMENT",
					`EnumField "EnumFieldName" should have a comment of the form "// EnumFieldName ..."`,
				),
			},
		},
		{
			name: "failures for proto with invalid fields without Golang style comments due to inline",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumBody: []parser.Visitee{
							&parser.EnumField{
								Ident: "EnumFieldName",
								InlineComment: &parser.Comment{
									Raw: "// EnumFieldName is special.",
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
					"ENUM_FIELDS_HAVE_COMMENT",
					`EnumField "EnumFieldName" should have a comment of the form "// EnumFieldName ..."`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewEnumFieldsHaveCommentRule(rule.Severity_Error, test.inputShouldFollowGolangStyle)

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
