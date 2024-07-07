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
								MessageName: "Inner_MessageName",
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
					"MESSAGE_NAMES_UPPER_CAMEL_CASE",
					string(rule.SeverityError),
					`Message name "messageName" must be UpperCamelCase like "MessageName"`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   150,
						Line:     7,
						Column:   15,
					},
					"MESSAGE_NAMES_UPPER_CAMEL_CASE",
					string(rule.SeverityError),
					`Message name "Inner_MessageName" must be UpperCamelCase like "InnerMessageName"`,
				),
				report.Failuref(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"MESSAGE_NAMES_UPPER_CAMEL_CASE",
					string(rule.SeverityError),
					`Message name "Message_name" must be UpperCamelCase like "MessageName"`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewMessageNamesUpperCamelCaseRule(rule.SeverityError, false, autodisable.Noop)

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

func TestMessageNamesUpperCamelCaseRule_Apply_fix(t *testing.T) {
	tests := []struct {
		name          string
		inputFilename string
		wantFilename  string
	}{
		{
			name:          "no fix for a correct proto",
			inputFilename: "upperCamelCase.proto",
			wantFilename:  "upperCamelCase.proto",
		},
		{
			name:          "fix for an incorrect proto",
			inputFilename: "invalid.proto",
			wantFilename:  "upperCamelCase.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := rules.NewMessageNamesUpperCamelCaseRule(rule.SeverityError, true, autodisable.Noop)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}

func TestMessageNamesUpperCamelCaseRule_Apply_disable(t *testing.T) {
	tests := []struct {
		name               string
		inputFilename      string
		inputPlacementType autodisable.PlacementType
		wantFilename       string
	}{
		{
			name:          "do nothing in case of no violations",
			inputFilename: "upperCamelCase.proto",
			wantFilename:  "upperCamelCase.proto",
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
			r := rules.NewMessageNamesUpperCamelCaseRule(rule.SeverityError, true, test.inputPlacementType)
			testApplyFix(t, r, test.inputFilename, test.wantFilename)
		})
	}
}
