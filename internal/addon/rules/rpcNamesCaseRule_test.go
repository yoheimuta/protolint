package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/protolint/internal/linter/config"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/report"
)

func TestRPCNamesCaseRule_Apply(t *testing.T) {
	tests := []struct {
		name            string
		inputProto      *parser.Proto
		inputConvention config.ConventionType
		wantFailures    []report.Failure
	}{
		{
			name: "no failures for proto without rpc",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
				},
			},
			inputConvention: config.ConventionLowerCamel,
		},
		{
			name: "no failures for proto with valid rpc",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "rpcName",
							},
						},
					},
				},
			},
			inputConvention: config.ConventionLowerCamel,
		},
		{
			name: "no failures for proto with valid rpc",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "RPC_NAME",
							},
						},
					},
				},
			},
			inputConvention: config.ConventionUpperSnake,
		},
		{
			name: "no failures for proto with valid rpc",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "rpc_name",
							},
						},
					},
				},
			},
			inputConvention: config.ConventionLowerSnake,
		},
		{
			name: "failures for proto with the option of LowerCamelCase",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "RPCName",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   100,
										Line:     5,
										Column:   10,
									},
								},
							},
							&parser.RPC{
								RPCName: "RPC_NAME",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   200,
										Line:     10,
										Column:   20,
									},
								},
							},
							&parser.RPC{
								RPCName: "rpc_name",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   300,
										Line:     20,
										Column:   30,
									},
								},
							},
						},
					},
				},
			},
			inputConvention: config.ConventionLowerCamel,
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"RPC_NAMES_CASE",
					`RPC name "RPCName" must be LowerCamelCase`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"RPC_NAMES_CASE",
					`RPC name "RPC_NAME" must be LowerCamelCase`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   300,
						Line:     20,
						Column:   30,
					},
					"RPC_NAMES_CASE",
					`RPC name "rpc_name" must be LowerCamelCase`,
				),
			},
		},
		{
			name: "failures for proto with the option of UpperSnakeCase",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "RPCName",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   100,
										Line:     5,
										Column:   10,
									},
								},
							},
							&parser.RPC{
								RPCName: "rpcName",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   200,
										Line:     10,
										Column:   20,
									},
								},
							},
							&parser.RPC{
								RPCName: "rpc_name",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   300,
										Line:     20,
										Column:   30,
									},
								},
							},
						},
					},
				},
			},
			inputConvention: config.ConventionUpperSnake,
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"RPC_NAMES_CASE",
					`RPC name "RPCName" must be UpperSnakeCase`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"RPC_NAMES_CASE",
					`RPC name "rpcName" must be UpperSnakeCase`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   300,
						Line:     20,
						Column:   30,
					},
					"RPC_NAMES_CASE",
					`RPC name "rpc_name" must be UpperSnakeCase`,
				),
			},
		},
		{
			name: "failures for proto with the option of LowerSnakeCase",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{
						ServiceBody: []parser.Visitee{
							&parser.RPC{
								RPCName: "RPCName",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   100,
										Line:     5,
										Column:   10,
									},
								},
							},
							&parser.RPC{
								RPCName: "rpcName",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   200,
										Line:     10,
										Column:   20,
									},
								},
							},
							&parser.RPC{
								RPCName: "RPC_NAME",
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   300,
										Line:     20,
										Column:   30,
									},
								},
							},
						},
					},
				},
			},
			inputConvention: config.ConventionLowerSnake,
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"RPC_NAMES_CASE",
					`RPC name "RPCName" must be LowerSnakeCase`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"RPC_NAMES_CASE",
					`RPC name "rpcName" must be LowerSnakeCase`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   300,
						Line:     20,
						Column:   30,
					},
					"RPC_NAMES_CASE",
					`RPC name "RPC_NAME" must be LowerSnakeCase`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewRPCNamesCaseRule(test.inputConvention)

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
