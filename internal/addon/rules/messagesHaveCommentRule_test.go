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

func TestMessagesHaveCommentRule_Apply(t *testing.T) {
	tests := []struct {
		name                         string
		inputProto                   *parser.Proto
		inputShouldFollowGolangStyle bool
		wantFailures                 []report.Failure
	}{
		{
			name: "no failures for proto without message",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
				},
			},
		},
		{
			name: "no failures for proto including valid messages with comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageName: "MessageName",
						Comments: []*parser.Comment{
							{
								Raw: "// a message name.",
							},
						},
					},
					&parser.Message{
						MessageName: "MessageName2",
						Comments: []*parser.Comment{
							{
								Raw: "// MessageName2 is a message name.",
							},
						},
					},
					&parser.Message{
						MessageName: "MessageName3",
						InlineComment: &parser.Comment{
							Raw: "// a message name.",
						},
					},
					&parser.Message{
						MessageName: "MessageName4",
						InlineCommentBehindLeftCurly: &parser.Comment{
							Raw: "// a message name.",
						},
					},
				},
			},
		},
		{
			name: "no failures for proto including valid messages with Golang style comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageName: "MessageName",
						Comments: []*parser.Comment{
							{
								Raw: "// MessageName is a message name.",
							},
						},
					},
					&parser.Message{
						MessageName: "MessageName2",
						Comments: []*parser.Comment{
							{
								Raw: "// MessageName2 is a message name.",
							},
						},
					},
				},
			},
			inputShouldFollowGolangStyle: true,
		},
		{
			name: "failures for proto with invalid messages",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageName: "MessageName",
						MessageBody: []parser.Visitee{
							&parser.Message{
								MessageName: "MessageName2",
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
					"MESSAGES_HAVE_COMMENT",
					string(rule.SeverityError),
					`Message "MessageName" should have a comment`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"MESSAGES_HAVE_COMMENT",
					string(rule.SeverityError),
					`Message "MessageName2" should have a comment`,
				),
			},
		},
		{
			name: "failures for proto with invalid messages without Golang style comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageName: "MessageName",
						MessageBody: []parser.Visitee{
							&parser.Message{
								MessageName: "MessageName2",
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
						Comments: []*parser.Comment{
							{
								Raw: "// a message name.",
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
					"MESSAGES_HAVE_COMMENT",
					string(rule.SeverityError),
					`Message "MessageName" should have a comment of the form "// MessageName ..."`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"MESSAGES_HAVE_COMMENT",
					string(rule.SeverityError),
					`Message "MessageName2" should have a comment of the form "// MessageName2 ..."`,
				),
			},
		},
		{
			name: "failures for proto with messages without Golang style comments due to the inline comment",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageName: "MessageName",
						MessageBody: []parser.Visitee{
							&parser.Message{
								MessageName: "MessageName2",
								InlineCommentBehindLeftCurly: &parser.Comment{
									Raw: "// MessageName2 is special.",
								},
							},
						},
						InlineComment: &parser.Comment{
							Raw: "// MessageName is special.",
						},
					},
				},
			},
			inputShouldFollowGolangStyle: true,
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{},
					"MESSAGES_HAVE_COMMENT",
					string(rule.SeverityError),
					`Message "MessageName" should have a comment of the form "// MessageName ..."`,
				),
				report.Failuref(
					meta.Position{},
					"MESSAGES_HAVE_COMMENT",
					string(rule.SeverityError),
					`Message "MessageName2" should have a comment of the form "// MessageName2 ..."`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewMessagesHaveCommentRule(rule.SeverityError, test.inputShouldFollowGolangStyle)

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
