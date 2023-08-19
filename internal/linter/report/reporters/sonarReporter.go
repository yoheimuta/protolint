package reporters

import (
	"encoding/json"
	"io"

	"github.com/yoheimuta/protolint/linter/report"
)

const (
	protolintSonarEngineId  string = "protolint"
	protolintSonarIssueType string = "CODE_SMELL"
	severityError           string = "MAJOR"
	severityWarn            string = "MINOR"
	severityNote            string = "INFO"
)

type SonarReporter struct {
}

// for details refer to https://docs.sonarsource.com/sonarqube/latest/analyzing-source-code/importing-external-issues/generic-issue-import-format/

type sonarTextRange struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
}

type sonarLocation struct {
	Message   string         `json:"message"`
	FilePath  string         `json:"filePath"`
	TextRange sonarTextRange `json:"textRange"`
}

type sonarIssue struct {
	EngineId        string        `json:"engineId"`
	RuleId          string        `json:"ruleId"`
	PrimaryLocation sonarLocation `json:"primaryLocation"`
	Severity        string        `json:"severity"`
	IssueType       string        `json:"issueType"`
}

func (s SonarReporter) Report(w io.Writer, fs []report.Failure) error {
	var issues []sonarIssue
	for _, f := range fs {
		issue := sonarIssue{
			EngineId:  protolintSonarEngineId,
			RuleId:    f.RuleID(),
			IssueType: protolintSonarIssueType,
			Severity:  getSonarSeverity(f.Severity()),
			PrimaryLocation: sonarLocation{
				Message:  f.Message(),
				FilePath: f.Pos().Filename,
				TextRange: sonarTextRange{
					StartLine:   f.Pos().Line,
					StartColumn: f.Pos().Column,
				},
			},
		}

		issues = append(issues, issue)
	}

	bs, err := json.MarshalIndent(issues, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write(bs)
	if err != nil {
		return err
	}

	return nil
}

func getSonarSeverity(severity string) string {
	if severity == "warn" {
		return severityWarn
	}

	if severity == "note" {
		return severityNote
	}

	return severityError
}
