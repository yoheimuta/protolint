package reporters_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/yoheimuta/go-protoparser/v4/parser/meta"
	"github.com/yoheimuta/protolint/internal/linter/report/reporters"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
)

type testFiles struct {
	files []string
}

type testStruct struct {
	name          string
	inputFailures []report.Failure
	wantOutput    string
	expectedError error
	files         testFiles
}

type testCases struct {
	oneWarningOneError testStruct
	onlyErrors         testStruct
	oneOfEach          testStruct
}

func makeTestData() testCases {
	return testCases{

		oneWarningOneError: testStruct{
			name:          "oneWarningOneError",
			wantOutput:    "",
			expectedError: nil,
			files:         testFiles{},
			inputFailures: []report.Failure{
				report.FailureWithSeverityf(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					string(rule.SeverityWarning),
					`EnumField name "fIRST_VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
				report.FailureWithSeverityf(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					string(rule.SeverityError),
					`EnumField name "SECOND.VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
			},
		},

		onlyErrors: testStruct{
			name:          "onlyErrors",
			wantOutput:    "",
			expectedError: nil,
			files:         testFiles{},
			inputFailures: []report.Failure{
				report.FailureWithSeverityf(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					string(rule.SeverityError),
					`EnumField name "fIRST_VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
				report.FailureWithSeverityf(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					string(rule.SeverityError),
					`EnumField name "SECOND.VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
				report.FailureWithSeverityf(
					meta.Position{
						Filename: "example.proto",
						Offset:   300,
						Line:     20,
						Column:   40,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					string(rule.SeverityError),
					`EnumField name "third.VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
			},
		},

		oneOfEach: testStruct{
			name:          "oneOfEach",
			wantOutput:    "",
			expectedError: nil,
			files:         testFiles{},
			inputFailures: []report.Failure{
				report.FailureWithSeverityf(
					meta.Position{
						Filename: "example.proto",
						Offset:   100,
						Line:     5,
						Column:   10,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					string(rule.SeverityNote),
					`EnumField name "fIRST_VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
				report.FailureWithSeverityf(
					meta.Position{
						Filename: "example.proto",
						Offset:   200,
						Line:     10,
						Column:   20,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					string(rule.SeverityWarning),
					`EnumField name "SECOND.VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
				report.FailureWithSeverityf(
					meta.Position{
						Filename: "example.proto",
						Offset:   300,
						Line:     20,
						Column:   40,
					},
					"ENUM_NAMES_UPPER_CAMEL_CASE",
					string(rule.SeverityError),
					`EnumField name "third.VALUE" must be CAPITALS_WITH_UNDERSCORES`,
				),
			},
		},
	}
}

func TestProblemMatcherReporter_Report(t *testing.T) {
	initTestCases := makeTestData()
	initTestCases.oneOfEach.want(`Protolint ENUM_NAMES_UPPER_CAMEL_CASE (info): example.proto[5,10]: EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (warning): example.proto[10,20]: EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[20,40]: EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.oneWarningOneError.want(`Protolint ENUM_NAMES_UPPER_CAMEL_CASE (warning): example.proto[5,10]: EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[10,20]: EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.onlyErrors.want(`Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[5,10]: EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[10,20]: EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[20,40]: EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)

	tests := initTestCases.tests()

	reporter := reporters.NewCiReporterWithGenericFormat()
	run_tests(t, tests, reporter)
}

func TestAzureDevOpsMatcherReporter_Report(t *testing.T) {
	initTestCases := makeTestData()
	initTestCases.oneOfEach.want(`##vso[task.logissue type=warning;sourcepath=example.proto;linenumber=10;columnnumber=20;code=ENUM_NAMES_UPPER_CAMEL_CASE;]EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
##vso[task.logissue type=error;sourcepath=example.proto;linenumber=20;columnnumber=40;code=ENUM_NAMES_UPPER_CAMEL_CASE;]EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.oneWarningOneError.want(`##vso[task.logissue type=warning;sourcepath=example.proto;linenumber=5;columnnumber=10;code=ENUM_NAMES_UPPER_CAMEL_CASE;]EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
##vso[task.logissue type=error;sourcepath=example.proto;linenumber=10;columnnumber=20;code=ENUM_NAMES_UPPER_CAMEL_CASE;]EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.onlyErrors.want(`##vso[task.logissue type=error;sourcepath=example.proto;linenumber=5;columnnumber=10;code=ENUM_NAMES_UPPER_CAMEL_CASE;]EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
##vso[task.logissue type=error;sourcepath=example.proto;linenumber=10;columnnumber=20;code=ENUM_NAMES_UPPER_CAMEL_CASE;]EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
##vso[task.logissue type=error;sourcepath=example.proto;linenumber=20;columnnumber=40;code=ENUM_NAMES_UPPER_CAMEL_CASE;]EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)

	tests := initTestCases.tests()

	reporter := reporters.NewCiReporterForAzureDevOps()
	run_tests(t, tests, reporter)
}
func TestGithubActionMatcherReporter_Report(t *testing.T) {
	initTestCases := makeTestData()
	initTestCases.oneOfEach.want(`::notice file=example.proto,line=5,col=10,title=ENUM_NAMES_UPPER_CAMEL_CASE::EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
::warning file=example.proto,line=10,col=20,title=ENUM_NAMES_UPPER_CAMEL_CASE::EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
::error file=example.proto,line=20,col=40,title=ENUM_NAMES_UPPER_CAMEL_CASE::EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.oneWarningOneError.want(`::warning file=example.proto,line=5,col=10,title=ENUM_NAMES_UPPER_CAMEL_CASE::EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
::error file=example.proto,line=10,col=20,title=ENUM_NAMES_UPPER_CAMEL_CASE::EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.onlyErrors.want(`::error file=example.proto,line=5,col=10,title=ENUM_NAMES_UPPER_CAMEL_CASE::EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
::error file=example.proto,line=10,col=20,title=ENUM_NAMES_UPPER_CAMEL_CASE::EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
::error file=example.proto,line=20,col=40,title=ENUM_NAMES_UPPER_CAMEL_CASE::EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)

	tests := initTestCases.tests()

	reporter := reporters.NewCiReporterForGithubActions()
	run_tests(t, tests, reporter)
}
func TestGitlabCiCdMatcherReporter_Report(t *testing.T) {
	initTestCases := makeTestData()
	initTestCases.oneOfEach.want(`INFO: ENUM_NAMES_UPPER_CAMEL_CASE  example.proto(5,10) : EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
WARNING: ENUM_NAMES_UPPER_CAMEL_CASE  example.proto(10,20) : EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
ERROR: ENUM_NAMES_UPPER_CAMEL_CASE  example.proto(20,40) : EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.oneWarningOneError.want(`WARNING: ENUM_NAMES_UPPER_CAMEL_CASE  example.proto(5,10) : EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
ERROR: ENUM_NAMES_UPPER_CAMEL_CASE  example.proto(10,20) : EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.onlyErrors.want(`ERROR: ENUM_NAMES_UPPER_CAMEL_CASE  example.proto(5,10) : EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
ERROR: ENUM_NAMES_UPPER_CAMEL_CASE  example.proto(10,20) : EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
ERROR: ENUM_NAMES_UPPER_CAMEL_CASE  example.proto(20,40) : EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)

	tests := initTestCases.tests()
	reporter := reporters.NewCiReporterForGitlab()
	run_tests(t, tests, reporter)
}

func TestEnvMatcherReporterFromTemplateString_Report(t *testing.T) {
	t.Setenv("PROTOLINT_CIREPORTER_TEMPLATE_STRING", "{{ .Severity }}@{{ .File }}[{{ .Line }},{{ .Column }}] triggered rule {{ .Rule }} with message {{ .Message }}")
	initTestCases := makeTestData()
	initTestCases.oneOfEach.want(`info@example.proto[5,10] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
warning@example.proto[10,20] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
error@example.proto[20,40] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.oneWarningOneError.want(`warning@example.proto[5,10] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
error@example.proto[10,20] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.onlyErrors.want(`error@example.proto[5,10] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
error@example.proto[10,20] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
error@example.proto[20,40] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)

	tests := initTestCases.tests()
	reporter := reporters.NewCiReporterFromEnv()
	run_tests(t, tests, reporter)
}

func TestEnvMatcherReporterFromTemplateFile_Report(t *testing.T) {
	initTestCases := makeTestData()
	initTestCases.oneOfEach.files.appendFile(t, "0oneOfEach.template", "[fromfile:0oneOfEach.template] {{ .Severity }}@{{ .File }}[{{ .Line }},{{ .Column }}] triggered rule {{ .Rule }} with message {{ .Message }}")
	initTestCases.oneOfEach.want(`[fromfile:0oneOfEach.template] info@example.proto[5,10] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
[fromfile:0oneOfEach.template] warning@example.proto[10,20] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
[fromfile:0oneOfEach.template] error@example.proto[20,40] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.oneWarningOneError.files.appendFile(t, "0oneWarningOneError.template", "[fromfile:0oneWarningOneError.template] {{ .Severity }}@{{ .File }}[{{ .Line }},{{ .Column }}] triggered rule {{ .Rule }} with message {{ .Message }}")
	initTestCases.oneWarningOneError.want(`[fromfile:0oneWarningOneError.template] warning@example.proto[5,10] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
[fromfile:0oneWarningOneError.template] error@example.proto[10,20] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.onlyErrors.files.appendFile(t, "0onlyErrors.template", "[fromfile:0onlyErrors.template] {{ .Severity }}@{{ .File }}[{{ .Line }},{{ .Column }}] triggered rule {{ .Rule }} with message {{ .Message }}")
	initTestCases.onlyErrors.want(`[fromfile:0onlyErrors.template] error@example.proto[5,10] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
[fromfile:0onlyErrors.template] error@example.proto[10,20] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
[fromfile:0onlyErrors.template] error@example.proto[20,40] triggered rule ENUM_NAMES_UPPER_CAMEL_CASE with message EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)

	tests := initTestCases.tests()
	reporter := reporters.NewCiReporterFromEnv()
	run_tests(t, tests, reporter)
}

func TestEnvMatcherReporterFromNonExistingTemplateFile_Report(t *testing.T) {
	t.Setenv("PROTOLINT_CIREPORTER_TEMPLATE_FILE", filepath.Join(t.TempDir(), "does.not.exist"))
	initTestCases := makeTestData()
	initTestCases.oneOfEach.want(`Protolint ENUM_NAMES_UPPER_CAMEL_CASE (info): example.proto[5,10]: EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (warning): example.proto[10,20]: EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[20,40]: EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.oneWarningOneError.want(`Protolint ENUM_NAMES_UPPER_CAMEL_CASE (warning): example.proto[5,10]: EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[10,20]: EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.onlyErrors.want(`Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[5,10]: EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[10,20]: EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[20,40]: EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)

	tests := initTestCases.tests()
	reporter := reporters.NewCiReporterFromEnv()
	run_tests(t, tests, reporter)
}

func TestEnvMatcherReporterFromUnallowedTemplateFile_Report(t *testing.T) {
	tplFile := filepath.Join(t.TempDir(), "does.exist")
	t.Setenv("PROTOLINT_CIREPORTER_TEMPLATE_FILE", tplFile)
	t.Cleanup(func() {
		if _, err := os.Stat(tplFile); os.IsNotExist(err) {
			t.Logf("File %s already deleted", tplFile)
		}

		err := os.Remove(tplFile)
		if err != nil {
			t.Errorf("ERROR while cleaning up file %s: %v", tplFile, err)
		}
	})
	file, err := os.Create(tplFile)
	if err != nil {
		t.Errorf("Failed to create temp file: %v", err)
	}
	_, err = file.Write([]byte("[fromfile:does.exist] {{ .Severity }}@{{ .File }}[{{ .Line }},{{ .Column }}] triggered rule {{ .Rule }} with message {{ .Message }}"))
	if err != nil {
		t.Errorf("Failed to write temp file: %v", err)
	}
	err = file.Close()
	if err != nil {
		t.Errorf("Failed to create temp file (error while closing): %v", err)
	}
	err = os.Chmod(tplFile, 0200)
	if err != nil {
		t.Errorf("Failed to remove read permissions: %v", err)
	}

	initTestCases := makeTestData()
	initTestCases.oneOfEach.want(`Protolint ENUM_NAMES_UPPER_CAMEL_CASE (info): example.proto[5,10]: EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (warning): example.proto[10,20]: EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[20,40]: EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.oneWarningOneError.want(`Protolint ENUM_NAMES_UPPER_CAMEL_CASE (warning): example.proto[5,10]: EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[10,20]: EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)
	initTestCases.onlyErrors.want(`Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[5,10]: EnumField name \"fIRST_VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[10,20]: EnumField name \"SECOND.VALUE\" must be CAPITALS_WITH_UNDERSCORES
Protolint ENUM_NAMES_UPPER_CAMEL_CASE (error): example.proto[20,40]: EnumField name \"third.VALUE\" must be CAPITALS_WITH_UNDERSCORES
`)

	tests := initTestCases.tests()
	reporter := reporters.NewCiReporterFromEnv()
	run_tests(t, tests, reporter)
}

func TestEnvMatcherReporterFromErroneousTemplateFile_ReportTestEnvMatcherReporterFromErroneousTemplateString_Report(t *testing.T) {
	initTestCases := makeTestData()
	initTestCases.oneOfEach.files.appendFile(t, "1oneOfEach.template", "[fromfile:1oneOfEach.template] {{ .Severit }}@{{ .Fil }}[{{ .Lne }},{{ .Colum }}] triggered rule {{ .ule }} with message {{ .Msg }}")
	initTestCases.oneOfEach.expectError(`template: Failure:1:34: executing "Failure" at <.Severit>: can't evaluate field Severit in type reporters.ciReportedFailure`)

	initTestCases.oneWarningOneError.files.appendFile(t, "1oneWarningOneError.template", "[fromfile:1oneWarningOneError.template] {{ .Severit }}@{{ .Fil }}[{{ .Lne }},{{ .Colum }}] triggered rule {{ .ule }} with message {{ .Msg }}")
	initTestCases.oneWarningOneError.expectError(`template: Failure:1:43: executing "Failure" at <.Severit>: can't evaluate field Severit in type reporters.ciReportedFailure`)

	initTestCases.onlyErrors.files.appendFile(t, "1onlyErrors.template", "[fromfile:1onlyErrors.template] {{ .Severit }}@{{ .Fil }}[{{ .Lne }},{{ .Colum }}] triggered rule {{ .ule }} with message {{ .Msg }}")
	initTestCases.onlyErrors.expectError(`template: Failure:1:35: executing "Failure" at <.Severit>: can't evaluate field Severit in type reporters.ciReportedFailure`)

	tests := initTestCases.tests()
	reporter := reporters.NewCiReporterFromEnv()
	run_tests(t, tests, reporter)
}

func TestEnvMatcherReporterFromErroneousTemplateString_Report(t *testing.T) {
	t.Setenv("PROTOLINT_CIREPORTER_TEMPLATE_STRING", "{{ .Severit }}@{{ .Fil }}[{{ .Lne }},{{ .Colum }}] triggered rule {{ .ule }} with message {{ .Msg }}")
	initTestCases := makeTestData()
	initTestCases.oneOfEach.expectError(`template: Failure:1:3: executing "Failure" at <.Severit>: can't evaluate field Severit in type reporters.ciReportedFailure`)
	initTestCases.oneWarningOneError.expectError(`template: Failure:1:3: executing "Failure" at <.Severit>: can't evaluate field Severit in type reporters.ciReportedFailure`)
	initTestCases.onlyErrors.expectError(`template: Failure:1:3: executing "Failure" at <.Severit>: can't evaluate field Severit in type reporters.ciReportedFailure`)

	tests := initTestCases.tests()
	reporter := reporters.NewCiReporterFromEnv()
	run_tests(t, tests, reporter)
}

func run_tests(t *testing.T, tests []testStruct, r reporters.CiReporter) {
	if len(tests) == 0 {
		t.SkipNow()
	}
	for _, test := range tests {
		test := test
		t.Cleanup(func() { test.cleanUp(t) })
		t.Run(test.name, func(t *testing.T) {
			if len(test.files.files) > 0 {
				t.Setenv("PROTOLINT_CIREPORTER_TEMPLATE_FILE", test.files.files[0])
			}

			buf := &bytes.Buffer{}
			err := r.Report(buf, test.inputFailures)
			test.cleanUp(t)
			if err != nil {
				isExpectedError := false
				if test.expectedError != nil {
					isExpectedError = err.Error() == test.expectedError.Error()
				}

				if !isExpectedError {
					t.Errorf("got err %v, but want %v", err, test.expectedError)
					return
				} else {
					t.Logf("Wanted %v, got %v", test.expectedError, err)
				}
			}
			if buf.String() != test.wantOutput {
				t.Errorf(`  got 
%s
, but want
%s`, buf.String(), test.wantOutput)
			}
		})
	}
}

func (ts testStruct) cleanUp(t *testing.T) {
	ts.files.cleanUp(t)
}

func (tf testFiles) cleanUp(t *testing.T) {
	for _, file := range tf.files {

		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Logf("File %s already deleted", file)
			continue
		}

		err := os.Remove(file)
		if err != nil {
			t.Errorf("ERROR while cleaning up file %s: %v", file, err)
		}
	}
}

func (tf *testFiles) appendFile(t *testing.T, fileName string, fileContent string) {
	filePath := filepath.Join(t.TempDir(), fileName)
	file, err := os.Create(filePath)
	if err != nil {
		t.Errorf("Failed to create file %s: %v", filePath, err)
		return
	}

	defer func(td *testing.T) {
		err = file.Close()
		if err != nil {
			td.Errorf("Failed to content to file %s: %v", filePath, err)
			return
		}
	}(t)

	_, err = file.WriteString(fileContent)
	if err != nil {
		t.Errorf("Failed to content to file %s: %v", filePath, err)
		return
	}

	tf.files = append(tf.files, filePath)
}

func (ts *testStruct) want(input string) {
	ts.wantOutput = input
}

func (ts *testStruct) expectError(err string) {
	ts.expectedError = fmt.Errorf(err)
	ts.want("")
}

func (td testCases) tests() []testStruct {
	var allCases []testStruct

	allCases = append(allCases, td.oneOfEach)
	allCases = append(allCases, td.oneWarningOneError)
	allCases = append(allCases, td.onlyErrors)

	return allCases
}
