package rules_test

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/protolint/internal/osutil"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/internal/linter/file"

	"github.com/yoheimuta/protolint/internal/setting_test"

	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/report"
)

func TestIndentRule_Apply(t *testing.T) {
	defaultSpace := strings.Repeat(" ", 2)

	tests := []struct {
		name           string
		inputStyle     string
		inputProtoPath string
		wantFailures   []report.Failure
		wantExistErr   bool
	}{
		{
			name:           "correct syntax",
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "syntax.proto"),
		},
		{
			name:           "incorrect syntax",
			inputStyle:     defaultSpace,
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "incorrect_syntax.proto"),
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_syntax.proto"),
						Offset:   14,
						Line:     2,
						Column:   5,
					},
					"INDENT",
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					"    ",
					"",
				),
			},
		},
		{
			name:           "correct enum",
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "enum.proto"),
		},
		{
			name:           "incorrect enum",
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "incorrect_enum.proto"),
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_enum.proto"),
						Offset:   162,
						Line:     7,
						Column:   2,
					},
					"INDENT",
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					" ",
					"",
				),
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_enum.proto"),
						Offset:   67,
						Line:     4,
						Column:   9,
					},
					"INDENT",
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					"        ",
					defaultSpace,
				),
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_enum.proto"),
						Offset:   114,
						Line:     6,
						Column:   6,
					},
					"INDENT",
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					"     ",
					defaultSpace,
				),
			},
		},
		{
			name:           "correct message",
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "message.proto"),
		},
		{
			name:           "incorrect message",
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "incorrect_message.proto"),
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_message.proto"),
						Offset:   100,
						Line:     6,
						Column:   3,
					},
					"INDENT",
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					"  ",
					strings.Repeat(defaultSpace, 2),
				),
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_message.proto"),
						Offset:   156,
						Line:     9,
						Column:   1,
					},
					"INDENT",
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					"",
					defaultSpace,
				),
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_message.proto"),
						Offset:   287,
						Line:     14,
						Column:   7,
					},
					"INDENT",
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					"      ",
					strings.Repeat(defaultSpace, 2),
				),
			},
		},
		{
			name:           "handle the proto containing extend. Fix https://github.com/yoheimuta/protolint/issues/63",
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "issue_63.proto"),
		},
		{
			name: `handle the case that the last rpc method of a service is having a statement block.
Fix https://github.com/yoheimuta/protolint/issues/74`,
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "issue_74.proto"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewIndentRule(
				test.inputStyle,
				"\n",
				false,
			)

			proto, err := file.NewProtoFile(test.inputProtoPath, test.inputProtoPath).Parse(false)
			if err != nil {
				t.Errorf(err.Error())
				return
			}

			got, err := rule.Apply(proto)
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

type testData struct {
	filePath   string
	originData []byte
}

func newTestIndentData(
	fileName string,
) (testData, error) {
	return newTestData(setting_test.TestDataPath("rules", "indentrule", fileName))
}

func newTestData(
	filePath string,
) (testData, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return testData{}, nil
	}
	return testData{
		filePath:   filePath,
		originData: data,
	}, nil
}

func (d testData) data() ([]byte, error) {
	return ioutil.ReadFile(d.filePath)
}

func (d testData) restore() error {
	newlineChar := "\n"
	lines := strings.Split(string(d.originData), newlineChar)
	return osutil.WriteLinesToExistingFile(d.filePath, lines, newlineChar)
}

func TestIndentRule_Apply_fix(t *testing.T) {
	space2 := strings.Repeat(" ", 2)

	correctSyntaxPath, err := newTestIndentData("syntax.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	incorrectSyntaxPath, err := newTestIndentData("incorrect_syntax.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	correctEnumPath, err := newTestIndentData("enum.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	incorrectEnumPath, err := newTestIndentData("incorrect_enum.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	correctMessagePath, err := newTestIndentData("message.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	incorrectMessagePath, err := newTestIndentData("incorrect_message.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	correctIssue99Path, err := newTestIndentData("issue_99.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	incorrectIssue99Path, err := newTestIndentData("incorrect_issue_99.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	incorrectIssue139Path, err := newTestIndentData("incorrect_issue_139.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	correctIssue139Path, err := newTestIndentData("correct_issue_139.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	tests := []struct {
		name            string
		inputTestData   testData
		wantCorrectData testData
	}{
		{
			name:            "correct syntax",
			inputTestData:   correctSyntaxPath,
			wantCorrectData: correctSyntaxPath,
		},
		{
			name:            "incorrect syntax",
			inputTestData:   incorrectSyntaxPath,
			wantCorrectData: correctSyntaxPath,
		},
		{
			name:            "correct enum",
			inputTestData:   correctEnumPath,
			wantCorrectData: correctEnumPath,
		},
		{
			name:            "incorrect enum",
			inputTestData:   incorrectEnumPath,
			wantCorrectData: correctEnumPath,
		},
		{
			name:            "correct message",
			inputTestData:   correctMessagePath,
			wantCorrectData: correctMessagePath,
		},
		{
			name:            "incorrect message",
			inputTestData:   incorrectMessagePath,
			wantCorrectData: correctMessagePath,
		},
		{
			name:            "correct issue_99",
			inputTestData:   correctIssue99Path,
			wantCorrectData: correctIssue99Path,
		},
		{
			name:            "incorrect issue_99",
			inputTestData:   incorrectIssue99Path,
			wantCorrectData: correctIssue99Path,
		},
		{
			name:            "do nothing against inner elements on the same line. Fix https://github.com/yoheimuta/protolint/issues/139",
			inputTestData:   incorrectIssue139Path,
			wantCorrectData: correctIssue139Path,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewIndentRule(
				space2,
				"\n",
				true,
			)

			proto, err := file.NewProtoFile(test.inputTestData.filePath, test.inputTestData.filePath).Parse(false)
			if err != nil {
				t.Errorf(err.Error())
				return
			}

			_, err = rule.Apply(proto)
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}

			got, err := test.inputTestData.data()
			if !reflect.DeepEqual(got, test.wantCorrectData.originData) {
				t.Errorf(
					"got %s(%v), but want %s(%v)",
					string(got), got,
					string(test.wantCorrectData.originData), test.wantCorrectData.originData,
				)
			}

			err = test.inputTestData.restore()
			if err != nil {
				t.Errorf("got err %v", err)
			}
		})
	}
}
