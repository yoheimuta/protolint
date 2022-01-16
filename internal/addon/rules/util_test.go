package rules_test

import (
	"github.com/yoheimuta/protolint/internal/linter/file"
	"github.com/yoheimuta/protolint/internal/setting_test"
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

	input, err := newTestData(setting_test.TestDataPath("rules", dataDir, inputFilename))
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	want, err := newTestData(setting_test.TestDataPath("rules", dataDir, wantFilename))
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	proto, err := file.NewProtoFile(input.filePath, input.filePath).Parse(false)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	_, err = r.Apply(proto)
	if err != nil {
		t.Errorf("got err %v, but want nil", err)
		return
	}

	got, err := input.data()
	if !reflect.DeepEqual(got, want.originData) {
		t.Errorf(
			"got %s(%v), but want %s(%v)",
			string(got), got,
			string(want.originData), want.originData,
		)
	}

	err = input.restore()
	if err != nil {
		t.Errorf("got err %v", err)
	}
}
