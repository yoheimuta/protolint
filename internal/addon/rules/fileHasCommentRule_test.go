package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/report"
)

func TestFileHasCommentRule_Apply(t *testing.T) {
	tests := []struct {
		name         string
		inputProto   *parser.Proto
		wantFailures []report.Failure
	}{
		{
			name: "no failures for proto starting with a doc comment",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Syntax{
						Comments: []*parser.Comment{
							{
								Raw: "// this file defines something-something.",
							},
						},
					},
				},
			},
		},
		{
			name: "a failure for proto without any doc comments",
			inputProto: &parser.Proto{
				ProtoBody: []parser.Visitee{
					&parser.Syntax{
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
					"FILE_HAS_COMMENT",
					`File should start with a doc comment`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewFileHasCommentRule()

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
