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

func TestEnumsHaveCommentRule_Apply(t *testing.T) {
	tests := []struct {
		name                         string
		inputProto                   *parser.Proto
		inputShouldFollowGolangStyle bool
		wantFailures                 []report.Failure
	}{
		{
			name: "no failures for proto without enum",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{},
			},
		},
		{
			name: "no failures for proto including valid enums with comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumName: "EnumName",
						Comments: []*parser.Comment{
							{
								Raw: "// a enum name.",
							},
						},
					},
					&parser.Enum{
						EnumName: "EnumName2",
						InlineComment: &parser.Comment{
							Raw: "// a enum name.",
						},
					},
					&parser.Enum{
						EnumName: "EnumName3",
						InlineCommentBehindLeftCurly: &parser.Comment{
							Raw: "// a enum name.",
						},
					},
				},
			},
		},
		{
			name: "no failures for proto including valid enums with Golang style comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumName: "EnumName",
						Comments: []*parser.Comment{
							{
								Raw: "// EnumName is a enum name.",
							},
						},
					},
				},
			},
			inputShouldFollowGolangStyle: true,
		},
		{
			name: "failures for proto with invalid enums",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumName: "EnumName",
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
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   150,
						Line:     7,
						Column:   15,
					},
					"ENUMS_HAVE_COMMENT",
					string(rule.SeverityError),
					`Enum "EnumName" should have a comment`,
				),
			},
		},
		{
			name: "failures for proto with invalid enums without Golang style comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumName: "EnumName",
						Comments: []*parser.Comment{
							{
								Raw: "// a enum name.",
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
			inputShouldFollowGolangStyle: true,
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   150,
						Line:     7,
						Column:   15,
					},
					"ENUMS_HAVE_COMMENT",
					string(rule.SeverityError),
					`Enum "EnumName" should have a comment of the form "// EnumName ..."`,
				),
			},
		},
		{
			name: "failures for proto with invalid enums without Golang style comments due to inline",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Enum{
						EnumName: "EnumName",
						InlineComment: &parser.Comment{
							Raw: "// EnumName is special.",
						},
					},
					&parser.Enum{
						EnumName: "EnumName2",
						InlineCommentBehindLeftCurly: &parser.Comment{
							Raw: "// EnumName2 is special.",
						},
					},
				},
			},
			inputShouldFollowGolangStyle: true,
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{},
					"ENUMS_HAVE_COMMENT",
					string(rule.SeverityError),
					`Enum "EnumName" should have a comment of the form "// EnumName ..."`,
				),
				report.Failuref(
					meta.Position{},
					"ENUMS_HAVE_COMMENT",
					string(rule.SeverityError),
					`Enum "EnumName2" should have a comment of the form "// EnumName2 ..."`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewEnumsHaveCommentRule(rule.SeverityError, test.inputShouldFollowGolangStyle)

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
