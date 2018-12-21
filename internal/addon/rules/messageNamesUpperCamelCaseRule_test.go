package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/go-protoparser/parser/meta"
	"github.com/yoheimuta/protolinter/internal/addon/rules"
	"github.com/yoheimuta/protolinter/internal/linter/report"
)

func TestMessageNamesUpperCamelCaseRule_Apply(t *testing.T) {
	tests := []struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
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
			name: "no failures for proto with valid message names",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageName: "MessageName",
						MessageBody: []parser.Visitee{
							&parser.Message{
								MessageName: "InnerMessageName",
							},
						},
					},
					&parser.Message{
						MessageName: "MessageName2",
					},
				},
			},
		},
		{
			name: "failures for proto with invalid message names",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Service{},
					&parser.Message{
						MessageName: "messageName",
						MessageBody: []parser.Visitee{
							&parser.Message{
								MessageName: "Inner.MessageName",
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
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   100,
								Line:     5,
								Column:   10,
							},
						},
					},
					&parser.Message{
						MessageName: "Message_name",
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
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					`Message name "messageName" must be UpperCamelCase`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   150,
						Line:     7,
						Column:   15,
					},
					`Message name "Inner.MessageName" must be UpperCamelCase`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					`Message name "Message_name" must be UpperCamelCase`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewMessageNamesUpperCamelCaseRule()

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
