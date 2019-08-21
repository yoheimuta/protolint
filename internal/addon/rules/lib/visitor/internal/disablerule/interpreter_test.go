package disablerule_test

import (
	"testing"

	"github.com/yoheimuta/protolint/internal/addon/rules/lib/visitor/internal/disablerule"

	"github.com/yoheimuta/go-protoparser/parser"
)

func TestInterpreter_Interpret(t *testing.T) {
	type inOut struct {
		name                string
		inputComments       []*parser.Comment
		inputInlineComments []*parser.Comment
		wantIsDisabled      bool
	}
	tests := []struct {
		name        string
		inputRuleID string
		inOuts      []inOut
	}{
		{
			name:        "disable:next comments skip ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
			inputRuleID: "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
			inOuts: []inOut{
				{
					name: "rule is enabled when there are no comments",
				},
				{
					name: "rule is enabled when there are no correct disable:next comments",
					inputComments: []*parser.Comment{
						{
							Raw: `// disable:next ENUM_FIELD_NAMES_UPPER_SNAKE_CASE`,
						},
					},
				},
				{
					name: "rule is enabled when the ruleID does not match it",
					inputComments: []*parser.Comment{
						{
							Raw: `// protolint:disable:next ENUM_NAMES_UPPER_CAMEL_CASE`,
						},
					},
				},
				{
					name: "rule is disabled when there is a disable:next comment with a ruleID",
					inputComments: []*parser.Comment{
						{
							Raw: `// protolint:disable:next ENUM_FIELD_NAMES_UPPER_SNAKE_CASE`,
						},
					},
					wantIsDisabled: true,
				},
				{
					name: "rule is disabled when there is a disable:next c-style comment with a ruleID",
					inputComments: []*parser.Comment{
						{
							Raw: `/*
protolint:disable:next ENUM_FIELD_NAMES_UPPER_SNAKE_CASE
*/`,
						},
					},
					wantIsDisabled: true,
				},
				{
					name: "rule is disabled when there are disable:next comments with ruleIDs",
					inputComments: []*parser.Comment{
						{
							Raw: `// protolint:disable:next ENUM_FIELD_NAMES_UPPER_SNAKE_CASE`,
						},
						{
							Raw: `// protolint:disable:next ENUM_NAMES_UPPER_CAMEL_CASE`,
						},
					},
					wantIsDisabled: true,
				},
			},
		},
		{
			name:        "disable:next comments skip SERVICE_NAMES_UPPER_CAMEL_CASE",
			inputRuleID: "SERVICE_NAMES_UPPER_CAMEL_CASE",
			inOuts: []inOut{
				{
					name: "rule is enabled when the ruleID does not match it",
					inputComments: []*parser.Comment{
						{
							Raw: `// protolint:disable:next ENUM_NAMES_UPPER_CAMEL_CASE`,
						},
					},
				},
				{
					name: "rule is disabled when there is a disable:next comment with a ruleID",
					inputComments: []*parser.Comment{
						{
							Raw: `// protolint:disable:next SERVICE_NAMES_UPPER_CAMEL_CASE`,
						},
					},
					wantIsDisabled: true,
				},
			},
		},
		{
			name:        "disable:this comment skips SERVICE_NAMES_UPPER_CAMEL_CASE",
			inputRuleID: "SERVICE_NAMES_UPPER_CAMEL_CASE",
			inOuts: []inOut{
				{
					name: "rule is enabled when the ruleID does not match it",
					inputInlineComments: []*parser.Comment{
						{
							Raw: `// protolint:disable:this ENUM_NAMES_UPPER_CAMEL_CASE`,
						},
					},
				},
				{
					name: "rule is disabled when there is a disable:this comment with a ruleID",
					inputInlineComments: []*parser.Comment{
						{
							Raw: `// protolint:disable:this SERVICE_NAMES_UPPER_CAMEL_CASE`,
						},
					},
					wantIsDisabled: true,
				},
			},
		},
		{
			name:        "disable SERVICE_NAMES_UPPER_CAMEL_CASE",
			inputRuleID: "SERVICE_NAMES_UPPER_CAMEL_CASE",
			inOuts: []inOut{
				{
					name: "rule is not disabled when there is a disable comment with another ruleID",
					inputComments: []*parser.Comment{
						{
							Raw: `// protolint:disable ENUM_FIELD_NAMES_UPPER_SNAKE_CASE`,
						},
					},
				},
				{
					name: "rule is disabled when there is a disable comment with a ruleID",
					inputComments: []*parser.Comment{
						{
							Raw: `// protolint:disable SERVICE_NAMES_UPPER_CAMEL_CASE`,
						},
					},
					wantIsDisabled: true,
				},
				{
					name:           "rule is always disabled after a disable comment",
					wantIsDisabled: true,
				},
				{
					name: "rule is disabled when there is a disable:next comment with a ruleID",
					inputComments: []*parser.Comment{
						{
							Raw: `// protolint:disable:next SERVICE_NAMES_UPPER_CAMEL_CASE`,
						},
					},
					wantIsDisabled: true,
				},
				{
					name:           "rule is always disabled after a disable comment",
					wantIsDisabled: true,
				},
			},
		},
		{
			name:        "enable SERVICE_NAMES_UPPER_CAMEL_CASE",
			inputRuleID: "SERVICE_NAMES_UPPER_CAMEL_CASE",
			inOuts: []inOut{
				{
					name: "rule is disabled when there is a disable comment with a ruleID",
					inputComments: []*parser.Comment{
						{
							Raw: `// protolint:disable SERVICE_NAMES_UPPER_CAMEL_CASE`,
						},
					},
					wantIsDisabled: true,
				},
				{
					name: "rule is not enabled when there is an enable comment with another ruleID",
					inputComments: []*parser.Comment{
						{
							Raw: `// protolint:enable ENUM_FIELD_NAMES_UPPER_SNAKE_CASE`,
						},
					},
					wantIsDisabled: true,
				},
				{
					name: "rule is enabled when there is an enable comment with a same ruleID",
					inputComments: []*parser.Comment{
						{
							Raw: `// protolint:enable SERVICE_NAMES_UPPER_CAMEL_CASE`,
						},
					},
				},
				{
					name: "rule is always enabled after an enable comment",
				},
				{
					name: "rule is disabled when there is a disable:next comment with a ruleID",
					inputComments: []*parser.Comment{
						{
							Raw: `// protolint:disable:next SERVICE_NAMES_UPPER_CAMEL_CASE`,
						},
					},
					wantIsDisabled: true,
				},
				{
					name: "rule is always enabled after an enable comment and a disable:next comment",
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			interpreter := disablerule.NewInterpreter(test.inputRuleID)

			for _, expect := range test.inOuts {
				got := interpreter.Interpret(expect.inputComments, expect.inputInlineComments...)
				if got != expect.wantIsDisabled {
					t.Errorf("[%s] got %v, but want %v", expect.name, got, expect.wantIsDisabled)
				}
			}
		})
	}
}
