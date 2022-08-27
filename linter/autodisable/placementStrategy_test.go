package autodisable_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
	"github.com/yoheimuta/protolint/internal/setting_test"
	"github.com/yoheimuta/protolint/internal/util_test"
	"github.com/yoheimuta/protolint/linter/autodisable"
)

type inputDisable struct {
	inputOffset   int
	inputComments []*parser.Comment
	inputInline   *parser.Comment
}

func TestPlacementStrategy_Disable(t *testing.T) {
	tests := []struct {
		name               string
		inputPlacementType autodisable.PlacementType
		inputFilename      string
		inputRoleID        string
		inputDisable       []inputDisable
		wantFilename       string
	}{
		{
			name:               "no auto disable",
			inputPlacementType: autodisable.Noop,
			inputRoleID:        "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
			inputDisable: []inputDisable{
				{
					inputOffset: 34,
				},
			},
			inputFilename: "invalid_enum_field_names.proto",
			wantFilename:  "invalid_enum_field_names.proto",
		},
		{
			name:               "add a new line comment",
			inputPlacementType: autodisable.Next,
			inputRoleID:        "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
			inputDisable: []inputDisable{
				{
					inputOffset: 34,
				},
			},
			inputFilename: "invalid_enum_field_names.proto",
			wantFilename:  "disabled_enum_field_names.proto",
		},
		{
			name:               "add new three line comments",
			inputPlacementType: autodisable.Next,
			inputRoleID:        "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
			inputDisable: []inputDisable{
				{
					inputOffset: 63,
				},
				{
					inputOffset: 115,
				},
				{
					inputOffset: 145,
				},
			},
			inputFilename: "invalid_many_enum_field_names.proto",
			wantFilename:  "disabled_many_enum_field_names.proto",
		},
		{
			name:               "not merge the comment",
			inputPlacementType: autodisable.Next,
			inputRoleID:        "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
			inputDisable: []inputDisable{
				{
					inputOffset: 99,
				},
				{
					inputOffset: 119,
				},
			},
			inputFilename: "invalid_enum_field_names_comment.proto",
			wantFilename:  "disabled_nomerge_enum_field_names.proto",
		},
		{
			name:               "add an inline comment",
			inputPlacementType: autodisable.ThisThenNext,
			inputRoleID:        "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
			inputDisable: []inputDisable{
				{
					inputOffset: 34,
				},
			},
			inputFilename: "invalid_enum_field_names.proto",
			wantFilename:  "disabled_inline_enum_field_names.proto",
		},
		{
			name:               "merge an inline comment",
			inputPlacementType: autodisable.ThisThenNext,
			inputRoleID:        "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
			inputDisable: []inputDisable{
				{
					inputOffset: 99,
				},
				{
					inputOffset: 119,
					inputInline: &parser.Comment{
						Raw:  `// protolint:disable:this ENUM_FIELD_NAMES_PREFIX`,
						Meta: meta.Meta{Pos: meta.Position{Offset: 142}},
					},
				},
			},
			inputFilename: "invalid_inline_disable_enum_field_names.proto",
			wantFilename:  "disabled_merge_inline_enum_field_names.proto",
		},
		{
			name:               "add an inline comment and a line comment",
			inputPlacementType: autodisable.ThisThenNext,
			inputRoleID:        "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
			inputDisable: []inputDisable{
				{
					inputOffset: 99,
				},
				{
					inputOffset: 119,
					inputInline: &parser.Comment{
						Raw: `// See the reference page.`,
					},
				},
			},
			inputFilename: "invalid_inline_enum_field_names.proto",
			wantFilename:  "disabled_inline_line_enum_field_names.proto",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			inputFilePath := setting_test.TestDataPath("autodisable", test.inputFilename)
			wantFilePath := setting_test.TestDataPath("autodisable", test.wantFilename)

			strategy, err := autodisable.NewPlacementStrategy(
				test.inputPlacementType,
				inputFilePath,
				test.inputRoleID,
			)
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}
			testDisable(t, strategy, test.inputDisable, inputFilePath, wantFilePath)
		})
	}
}

func testDisable(
	t *testing.T,
	strategy autodisable.PlacementStrategy,
	inputDisable []inputDisable,
	inputFilePath string,
	wantFilePath string,
) {
	input, err := util_test.NewTestData(inputFilePath)
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	want, err := util_test.NewTestData(wantFilePath)
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	for _, p := range inputDisable {
		strategy.Disable(p.inputOffset, p.inputComments, p.inputInline)
	}
	err = strategy.Finalize()
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	got, _ := input.Data()
	if !reflect.DeepEqual(got, want.OriginData) {
		t.Errorf(
			"got %s(%v), but want %s(%v)",
			string(got), got,
			string(want.OriginData), want.OriginData,
		)
	}

	err = input.Restore()
	if err != nil {
		t.Errorf("got err %v", err)
	}
}
