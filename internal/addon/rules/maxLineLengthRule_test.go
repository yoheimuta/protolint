package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/internal/setting_test"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/report"
)

func TestMaxLineLengthRule_Apply(t *testing.T) {
	tests := []struct {
		name          string
		inputMaxChars int
		inputTabChars int
		inputProto    *parser.Proto
		wantFailures  []report.Failure
		wantExistErr  bool
	}{
		{
			name: "not found proto file",
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: "",
				},
			},
			wantExistErr: true,
		},
		{
			name:          "not found long lines",
			inputMaxChars: 120,
			inputTabChars: 4,
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: setting_test.TestDataPath("rules", "max_line_length_rule.proto"),
				},
			},
		},
		{
			name:          "found long lines",
			inputTabChars: 4,
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: setting_test.TestDataPath("rules", "max_line_length_rule.proto"),
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "max_line_length_rule.proto"),
						Line:     3,
						Column:   1,
					},
					"MAX_LINE_LENGTH",
					`The line length is 91, but it must be shorter than 80`,
				),
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "max_line_length_rule.proto"),
						Line:     15,
						Column:   1,
					},
					"MAX_LINE_LENGTH",
					`The line length is 88, but it must be shorter than 80`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewMaxLineLengthRule(
				test.inputMaxChars,
				test.inputTabChars,
			)

			got, err := rule.Apply(test.inputProto)
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
