package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/go-protoparser/parser/meta"
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/internal/linter/report"
)

func TestPackageNameLowerCaseRule_Apply(t *testing.T) {
	tests := []struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		{
			name: "no failures for proto without service",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Message{},
				},
			},
		},
		{
			name: "no failures for proto with the valid package name",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Package{
						Name: "package",
					},
				},
			},
		},
		{
			name: "no failures for proto with the valid package name with periods",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Package{
						Name: "my.v1.package",
					},
				},
			},
		},
		{
			name: "failures for proto with the invalid package name",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Package{
						Name: "myV1Package",
						Meta: meta.Meta{
							Pos: meta.Position{
								Filename: "example.proto",
								Offset:   100,
								Line:     5,
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
						Line:     5,
						Column:   10,
					},
					`Package name "myV1Package" must only contains lowercase letters, digits and/or periods.`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewPackageNameLowerCaseRule()

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
