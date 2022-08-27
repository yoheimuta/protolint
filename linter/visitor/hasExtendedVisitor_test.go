package visitor_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/internal/linter/file"
	"github.com/yoheimuta/protolint/internal/setting_test"
	"github.com/yoheimuta/protolint/internal/util_test"
	"github.com/yoheimuta/protolint/linter/autodisable"
	"github.com/yoheimuta/protolint/linter/visitor"

	"github.com/yoheimuta/protolint/linter/report"

	"github.com/yoheimuta/go-protoparser/v4/parser"
)

type testVisitor struct {
	*visitor.BaseAddVisitor
	next bool
}

func (v *testVisitor) VisitMessage(message *parser.Message) bool {
	v.AddFailuref(message.Meta.Pos, "Test Message")
	return v.next
}

type testVisitorInvalidEnumField struct {
	*visitor.BaseAddVisitor
	next bool
}

func (v *testVisitorInvalidEnumField) VisitEnumField(field *parser.EnumField) bool {
	v.AddFailuref(field.Meta.Pos, "Failed field")
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
				BaseAddVisitor: visitor.NewBaseAddVisitor("MESSAGE_NAMES_UPPER_CAMEL_CASE"),
			},
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{Filename: ""},
			},
		},
		{
			name: "visit a message",
			inputVisitor: &testVisitor{
				BaseAddVisitor: visitor.NewBaseAddVisitor("MESSAGE_NAMES_UPPER_CAMEL_CASE"),
			},
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{Filename: ""},
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
					"MESSAGE_NAMES_UPPER_CAMEL_CASE",
					"Test Message",
				),
			},
		},
		{
			name: "visit messages",
			inputVisitor: &testVisitor{
				BaseAddVisitor: visitor.NewBaseAddVisitor("MESSAGE_NAMES_UPPER_CAMEL_CASE"),
			},
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{Filename: ""},
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
					"MESSAGE_NAMES_UPPER_CAMEL_CASE",
					"Test Message",
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     20,
						Column:   10,
					},
					"MESSAGE_NAMES_UPPER_CAMEL_CASE",
					"Test Message",
				),
			},
		},
		{
			name: "visit messages recursively",
			inputVisitor: &testVisitor{
				BaseAddVisitor: visitor.NewBaseAddVisitor("MESSAGE_NAMES_UPPER_CAMEL_CASE"),
				next:           true,
			},
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{Filename: ""},
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
					"MESSAGE_NAMES_UPPER_CAMEL_CASE",
					"Test Message",
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     20,
						Column:   10,
					},
					"MESSAGE_NAMES_UPPER_CAMEL_CASE",
					"Test Message",
				),
			},
		},
		{
			name: "visit a message. one is disabled.",
			inputVisitor: &testVisitor{
				BaseAddVisitor: visitor.NewBaseAddVisitor("MESSAGE_NAMES_UPPER_CAMEL_CASE"),
			},
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{Filename: ""},
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
					"MESSAGE_NAMES_UPPER_CAMEL_CASE",
					"Test Message",
				),
			},
		},
		{
			name: "visit a message. one is disabled by an inline comment.",
			inputVisitor: &testVisitor{
				BaseAddVisitor: visitor.NewBaseAddVisitor("MESSAGE_NAMES_UPPER_CAMEL_CASE"),
			},
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{Filename: ""},
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
						InlineComment: &parser.Comment{
							Raw: `// protolint:disable:this MESSAGE_NAMES_UPPER_CAMEL_CASE`,
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
					"MESSAGE_NAMES_UPPER_CAMEL_CASE",
					"Test Message",
				),
			},
		},
		{
			name: "visit messages. others are disabled.",
			inputVisitor: &testVisitor{
				BaseAddVisitor: visitor.NewBaseAddVisitor("MESSAGE_NAMES_UPPER_CAMEL_CASE"),
			},
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{Filename: ""},
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
					"MESSAGE_NAMES_UPPER_CAMEL_CASE",
					"Test Message",
				),
			},
		},
		{
			name: "visit messages. others are disabled by a last line comment.",
			inputVisitor: &testVisitor{
				BaseAddVisitor: visitor.NewBaseAddVisitor("MESSAGE_NAMES_UPPER_CAMEL_CASE"),
			},
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{Filename: ""},
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
					&parser.Comment{
						Raw: `// protolint:enable MESSAGE_NAMES_UPPER_CAMEL_CASE`,
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
					"MESSAGE_NAMES_UPPER_CAMEL_CASE",
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

func TestRunVisitorAutoDisable(t *testing.T) {
	tests := []struct {
		name               string
		inputVisitor       visitor.HasExtendedVisitor
		inputFilename      string
		inputRuleID        string
		inputPlacementType autodisable.PlacementType
		wantExistErr       bool
		wantFailureCount   int
		wantFilename       string
	}{
		{
			name: "Do nothing in case of no failures",
			inputVisitor: &testVisitor{
				BaseAddVisitor: visitor.NewBaseAddVisitor("ENUM_FIELD_NAMES_UPPER_SNAKE_CASE"),
			},
			inputFilename:      "invalid.proto",
			inputRuleID:        "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
			inputPlacementType: autodisable.Next,
			wantFilename:       "invalid.proto",
		},
		{
			name: "Insert a disable:next comment",
			inputVisitor: &testVisitorInvalidEnumField{
				BaseAddVisitor: visitor.NewBaseAddVisitor("ENUM_FIELD_NAMES_UPPER_SNAKE_CASE"),
			},
			inputFilename:      "invalid.proto",
			inputRuleID:        "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
			inputPlacementType: autodisable.Next,
			wantFailureCount:   1,
			wantFilename:       "disable_next.proto",
		},
		{
			name: "Insert a disable:this comment",
			inputVisitor: &testVisitorInvalidEnumField{
				BaseAddVisitor: visitor.NewBaseAddVisitor("ENUM_FIELD_NAMES_UPPER_SNAKE_CASE"),
			},
			inputFilename:      "invalid.proto",
			inputRuleID:        "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
			inputPlacementType: autodisable.ThisThenNext,
			wantFailureCount:   1,
			wantFilename:       "disable_this.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			input, err := util_test.NewTestData(setting_test.TestDataPath("visitor", test.inputFilename))
			if err != nil {
				t.Errorf("got err %v", err)
				return
			}

			want, err := util_test.NewTestData(setting_test.TestDataPath("visitor", test.wantFilename))
			if err != nil {
				t.Errorf("got err %v", err)
				return
			}

			proto, err := file.NewProtoFile(input.FilePath, input.FilePath).Parse(false)
			if err != nil {
				t.Errorf(err.Error())
				return
			}

			got, err := visitor.RunVisitorAutoDisable(
				test.inputVisitor,
				proto,
				test.inputRuleID,
				test.inputPlacementType,
			)

			if test.wantExistErr {
				if err == nil {
					t.Errorf("got err nil, but want err")
				}
				return
			} else if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}

			if len(got) != test.wantFailureCount {
				t.Errorf("len(got) %v, but want %v", len(got), test.wantFailureCount)
			}

			got2, _ := input.Data()
			if !reflect.DeepEqual(got2, want.OriginData) {
				t.Errorf(
					"got %s(%v), but want %s(%v)",
					string(got2), got,
					string(want.OriginData), want.OriginData,
				)
			}

			err = input.Restore()
			if err != nil {
				t.Errorf("got err %v", err)
			}

		})
	}
}
