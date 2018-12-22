package visitor_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/parser/meta"

	"github.com/yoheimuta/protolint/internal/addon/rules/internal/visitor"

	"github.com/yoheimuta/protolint/internal/linter/report"

	"github.com/yoheimuta/go-protoparser/parser"
)

type testVisitor struct {
	*visitor.BaseAddVisitor
	next bool
}

func (v *testVisitor) VisitMessage(message *parser.Message) bool {
	v.AddFailuref(message.Meta.Pos, "Test Message")
	return v.next
}

func TestRunVisitor(t *testing.T) {
	tests := []struct {
		name         string
		inputVisitor *testVisitor
		inputProto   *parser.Proto
		inputRuleID  string
		wantExistErr bool
		wantFailures []report.Failure
	}{
		{
			name: "visit no messages",
			inputVisitor: &testVisitor{
				BaseAddVisitor: visitor.NewBaseAddVisitor(),
			},
			inputProto: &parser.Proto{},
		},
		{
			name: "visit a message",
			inputVisitor: &testVisitor{
				BaseAddVisitor: visitor.NewBaseAddVisitor(),
			},
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   100,
								Line:     10,
								Column:   5,
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
						Line:     10,
						Column:   5,
					},
					"Test Message",
				),
			},
		},
		{
			name: "visit messages",
			inputVisitor: &testVisitor{
				BaseAddVisitor: visitor.NewBaseAddVisitor(),
			},
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   100,
								Line:     10,
								Column:   5,
							},
						},
					},
					&parser.Message{
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   200,
								Line:     20,
								Column:   10,
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
						Line:     10,
						Column:   5,
					},
					"Test Message",
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     20,
						Column:   10,
					},
					"Test Message",
				),
			},
		},
		{
			name: "visit messages recursively",
			inputVisitor: &testVisitor{
				BaseAddVisitor: visitor.NewBaseAddVisitor(),
				next:           true,
			},
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						MessageBody: []parser.Visitee{
							&parser.Message{
								Meta: meta.Meta{
									Pos: meta.Position{
										Filename: "example.proto",
										Offset:   200,
										Line:     20,
										Column:   10,
									},
								},
							},
						},
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   100,
								Line:     10,
								Column:   5,
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
						Line:     10,
						Column:   5,
					},
					"Test Message",
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     20,
						Column:   10,
					},
					"Test Message",
				),
			},
		},
		{
			name: "visit a message. one is disabled.",
			inputVisitor: &testVisitor{
				BaseAddVisitor: visitor.NewBaseAddVisitor(),
			},
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   100,
								Line:     10,
								Column:   5,
							},
						},
						Comments: []*parser.Comment{
							{
								Raw: `// protolint:disable:next MESSAGE_NAMES_UPPER_CAMEL_CASE`,
							},
						},
					},
					&parser.Message{
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   200,
								Line:     20,
								Column:   10,
							},
						},
					},
				},
			},
			inputRuleID: `MESSAGE_NAMES_UPPER_CAMEL_CASE`,
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     20,
						Column:   10,
					},
					"Test Message",
				),
			},
		},
		{
			name: "visit a messages. others are disabled.",
			inputVisitor: &testVisitor{
				BaseAddVisitor: visitor.NewBaseAddVisitor(),
			},
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   100,
								Line:     10,
								Column:   5,
							},
						},
						Comments: []*parser.Comment{
							{
								Raw: `// protolint:disable MESSAGE_NAMES_UPPER_CAMEL_CASE`,
							},
						},
					},
					&parser.Message{
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   200,
								Line:     20,
								Column:   10,
							},
						},
					},
					&parser.Message{
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   300,
								Line:     30,
								Column:   15,
							},
						},
						Comments: []*parser.Comment{
							{
								Raw: `// protolint:enable MESSAGE_NAMES_UPPER_CAMEL_CASE`,
							},
						},
					},
				},
			},
			inputRuleID: `MESSAGE_NAMES_UPPER_CAMEL_CASE`,
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   300,
						Line:     30,
						Column:   15,
					},
					"Test Message",
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, err := visitor.RunVisitor(
				test.inputVisitor,
				test.inputProto,
				test.inputRuleID,
			)

			if test.wantExistErr {
				if err == nil {
					t.Errorf("got err nil, but want err")
				}
				return
			}
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
