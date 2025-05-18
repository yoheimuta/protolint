package rules_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/yoheimuta/protolint/internal/util_test"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/internal/linter/file"

	"github.com/yoheimuta/protolint/internal/setting_test"

	"github.com/yoheimuta/protolint/internal/addon/rules"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
)

func TestIndentRule_Apply(t *testing.T) {
	defaultSpace := strings.Repeat(" ", 2)

	tests := []struct {
		name               string
		inputStyle         string
		inputProtoPath     string
		inputInsertNewline bool
		wantFailures       []report.Failure
		wantExistErr       bool
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
					string(rule.SeverityError),
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
						Offset:   67,
						Line:     4,
						Column:   9,
					},
					"INDENT",
					string(rule.SeverityError),
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
					string(rule.SeverityError),
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					"     ",
					defaultSpace,
				),
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_enum.proto"),
						Offset:   162,
						Line:     7,
						Column:   2,
					},
					"INDENT",
					string(rule.SeverityError),
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					" ",
					"",
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
					string(rule.SeverityError),
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
					string(rule.SeverityError),
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
					string(rule.SeverityError),
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
		{
			name: `skip wrong indentations of inner elements on the same line.
Fix https://github.com/yoheimuta/protolint/issues/139`,
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "issue_139.proto"),
		},
		{
			name: `detect only a toplevel indentation mistake and skip other than that on the same line.
Fix https://github.com/yoheimuta/protolint/issues/139`,
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "incorrect_issue_139.proto"),
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_issue_139.proto"),
						Offset:   222,
						Line:     11,
						Column:   3,
					},
					"INDENT",
					string(rule.SeverityError),
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					"  ",
					"",
				),
			},
		},
		{
			name: `do not skip wrong indentations of inner elements on the same line.
Fix https://github.com/yoheimuta/protolint/issues/139`,
			inputProtoPath:     setting_test.TestDataPath("rules", "indentrule", "incorrect_issue_139_short.proto"),
			inputInsertNewline: true,
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_issue_139_short.proto"),
						Offset:   82,
						Line:     7,
						Column:   3,
					},
					"INDENT",
					string(rule.SeverityError),
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					"  ",
					"",
				),
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_issue_139_short.proto"),
						Offset:   104,
						Line:     7,
						Column:   25,
					},
					"INDENT",
					string(rule.SeverityError),
					`Found a possible incorrect indentation style. Inserting a new line is recommended.`,
				),
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "incorrect_issue_139_short.proto"),
						Offset:   127,
						Line:     7,
						Column:   48,
					},
					"INDENT",
					string(rule.SeverityError),
					`Found a possible incorrect indentation style. Inserting a new line is recommended.`,
				),
			},
		},
		{
			name: `handle the case that the proto has a mixture of line ending formats like LF and CRLF.
Fix https://github.com/yoheimuta/protolint/issues/280`,
			inputProtoPath: setting_test.TestDataPath("rules", "indentrule", "issue_280_mix_lineending.proto"),
			wantFailures: []report.Failure{
				report.Failuref(
					meta.Position{
						Filename: setting_test.TestDataPath("rules", "indentrule", "issue_280_mix_lineending.proto"),
						Offset:   580,
						Line:     27,
						Column:   5,
					},
					"INDENT",
					string(rule.SeverityError),
					`Found an incorrect indentation style "%s". "%s" is correct.`,
					"    ",
					"  ",
				),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule := rules.NewIndentRule(
				rule.SeverityError,
				test.inputStyle,
				!test.inputInsertNewline,
				false,
			)

			proto, err := file.NewProtoFile(test.inputProtoPath, test.inputProtoPath).Parse(false)
			if err != nil {
				t.Errorf("%v", err)
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
				if len(got) != len(test.wantFailures) {
					t.Errorf("len(got) %v, but len(want) %v", len(got), len(test.wantFailures))
					return
				}
				for k, v := range got {
					if !reflect.DeepEqual(v.Pos(), test.wantFailures[k].Pos()) {
						t.Errorf("got[%v].Pos() %v(offset=%v), but want[%v].Pos() %v", k, v.Pos(), v.Pos().Offset, k, test.wantFailures[k].Pos())
						continue
					}
					if !reflect.DeepEqual(v.Message(), test.wantFailures[k].Message()) {
						t.Errorf("got[%v].Message() %v, but want[%v].Message() %v", k, v.Message(), k, test.wantFailures[k].Message())
						continue
					}
					if !reflect.DeepEqual(v, test.wantFailures[k]) {
						t.Errorf("got[%v] %v, but want[%v] %v", k, v, k, test.wantFailures[k])
						continue
					}
				}
			}
		})
	}
}

func newTestIndentData(
	fileName string,
) (util_test.TestData, error) {
	return util_test.NewTestData(setting_test.TestDataPath("rules", "indentrule", fileName))
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

	correctIssue139Path, err := newTestIndentData("issue_139.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	correctIssue139InsertPath, err := newTestIndentData("issue_139_insert_linebreaks.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	issue409Path, err := newTestIndentData("issue_409.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	issue409FixedPath, err := newTestIndentData("issue_409_fixed.proto")
	if err != nil {
		t.Errorf("got err %v", err)
		return
	}

	tests := []struct {
		name               string
		inputTestData      util_test.TestData
		inputInsertNewline bool
		wantCorrectData    util_test.TestData
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
		{
			name:               "insert linebreaks against inner elements on the same line. Fix https://github.com/yoheimuta/protolint/issues/139",
			inputTestData:      incorrectIssue139Path,
			inputInsertNewline: true,
			wantCorrectData:    correctIssue139InsertPath,
		},
		{
			name:               "handle comments around repeated keyword without panic. Fix https://github.com/yoheimuta/protolint/issues/409",
			inputTestData:      issue409Path,
			inputInsertNewline: true,
			wantCorrectData:    issue409FixedPath,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rule_to_test := rules.NewIndentRule(
				rule.SeverityError,
				space2,
				!test.inputInsertNewline,
				true,
			)

			proto, err := file.NewProtoFile(test.inputTestData.FilePath, test.inputTestData.FilePath).Parse(false)
			if err != nil {
				t.Errorf("%v", err)
				return
			}

			_, err = rule_to_test.Apply(proto)
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}

			got, err := test.inputTestData.Data()
			if !reflect.DeepEqual(got, test.wantCorrectData.OriginData) {
				t.Errorf(
					"got %s(%v), but want %s(%v)",
					string(got), got,
					string(test.wantCorrectData.OriginData), test.wantCorrectData.OriginData,
				)
			}

			// restore the file
			defer func() {
				err = test.inputTestData.Restore()
				if err != nil {
					t.Errorf("got err %v", err)
				}
			}()

			// check whether the modified content can pass the lint in the end.
			ruleOnlyCheck := rules.NewIndentRule(
				rule.SeverityError,
				space2,
				!test.inputInsertNewline,
				false,
			)
			proto, err = file.NewProtoFile(test.inputTestData.FilePath, test.inputTestData.FilePath).Parse(false)
			if err != nil {
				t.Errorf("%v", err)
				return
			}
			gotCheck, err := ruleOnlyCheck.Apply(proto)
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}
			if 0 < len(gotCheck) {
				t.Errorf("got failures %v, but want no failures", gotCheck)
				return
			}
		})
	}
}
