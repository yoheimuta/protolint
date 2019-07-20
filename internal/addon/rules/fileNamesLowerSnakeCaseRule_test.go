package rules_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/parser/meta"

	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/internal/linter/report"
)

func TestFileNamesLowerSnakeCaseRule_Apply(t *testing.T) {
	tests := []struct {
		name          string
		inputProto    *parser.Proto
		inputExcluded []string
		wantFailures  []report.Failure
	}{
		{
			name: "no failures for proto with a valid file name",
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: "../proto/simple.proto",
				},
			},
		},
		{
			name: "no failures for proto with a valid lower snake case file name",
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: "../proto/lower_snake_case.proto",
				},
			},
		},
		{
			name: "no failures for excluded proto",
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: "proto/lowerSnakeCase.proto",
				},
			},
			inputExcluded: []string{
				"proto/lowerSnakeCase.proto",
			},
		},
		{
			name: "a failure for proto with a camel case file name",
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: "proto/lowerSnakeCase.proto",
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "proto/lowerSnakeCase.proto",
						Offset:   0,
						Line:     1,
						Column:   1,
					},
					`File name should be lower_snake_case.proto.`,
				),
			},
		},
		{
			name: "a failure for proto with an invalid file extension",
			inputProto: &parser.Proto{
				Meta: &parser.ProtoMeta{
					Filename: "proto/lowerSnakeCase.txt",
				},
			},
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: "proto/lowerSnakeCase.txt",
						Offset:   0,
						Line:     1,
						Column:   1,
					},
					`File name should be lower_snake_case.proto.`,
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewFileNamesLowerSnakeCaseRule(test.inputExcluded)

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
