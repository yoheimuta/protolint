package rules_test

import (
	"github.com/yoheimuta/protolint/internal/linter/file"
	"github.com/yoheimuta/protolint/internal/setting_test"
	"github.com/yoheimuta/protolint/internal/util_test"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/strs"

	"reflect"
	"testing"
)

func testApplyFix(
	t *testing.T,
	r rule.Rule,
	inputFilename string,
	wantFilename string,
) {
	dataDir := strs.ToLowerCamelCase(r.ID())

	input, err := util_test.NewTestData(setting_test.TestDataPath("rules", dataDir, inputFilename))
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	want, err := util_test.NewTestData(setting_test.TestDataPath("rules", dataDir, wantFilename))
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	proto, err := file.NewProtoFile(input.FilePath, input.FilePath).Parse(false)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	_, err = r.Apply(proto)
	if err != nil {
		t.Errorf("got err %v, but want nil", err)
		return
	}

	got, err := input.Data()
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
