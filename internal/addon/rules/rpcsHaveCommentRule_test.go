package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/report"
)

func TestRPCsHaveCommentRule_Apply(t *testing.T) {
	tests := []struct {
		name                         string
		inputProto                   *parser.Proto
		inputShouldFollowGolangStyle bool
		wantFailures                 []report.Failure
	}{
		{
			name: "no failures for proto without rpc",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{},
			},
		},
		{
			name: "no failures for proto including valid rpcs with comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "RPCName",
								Comments: []*parser.Comment{
									{
										Raw: "// a rpc name.",
									},
								},
							},
							&parser.RPC{
								RPCName: "RPCName2",
								InlineComment: &parser.Comment{
									Raw: "// a rpc name.",
								},
							},
							&parser.RPC{
								RPCName: "RPCName3",
								InlineCommentBehindLeftCurly: &parser.Comment{
									Raw: "// a rpc name.",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "no failures for proto including valid rpcs with Golang style comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "RPCName",
								Comments: []*parser.Comment{
									{
										Raw: "// RPCName is a rpc name.",
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
			name: "failures for proto with invalid rpcs",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "RPCName",
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
					"RPCS_HAVE_COMMENT",
					`RPC "RPCName" should have a comment`,
				),
			},
		},
		{
			name: "failures for proto with invalid rpcs without Golang style comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "RPCName",
								Comments: []*parser.Comment{
									{
										Raw: "// a rpc name.",
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
					"RPCS_HAVE_COMMENT",
					`RPC "RPCName" should have a comment of the form "// RPCName ..."`,
				),
			},
		},
		{
			name: "failures for proto with invalid rpcs without Golang style comments due to inline",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "RPCName",
								InlineComment: &parser.Comment{
									Raw: "// RPCName is special.",
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
					"RPCS_HAVE_COMMENT",
					`RPC "RPCName" should have a comment of the form "// RPCName ..."`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewRPCsHaveCommentRule(test.inputShouldFollowGolangStyle)

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
