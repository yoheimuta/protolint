package reporters

import (
	"io"

	"github.com/chavacava/garif"
	"github.com/yoheimuta/protolint/linter/report"
)

// SarifReporter creates reports formatted as a JSON
// Document.
// The document format is used according to the SARIF
// Standard.
// Refer to http://docs.oasis-open.org/sarif/sarif/v2.1.0/sarif-v2.1.0.html
// for details to the format.
type SarifReporter struct{}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Report writes failures to w formatted as a SARIF document.
func (r SarifReporter) Report(w io.Writer, fs []report.Failure) error {
	rulesByID := make(map[string]*garif.ReportingDescriptor)
	allRules := []*garif.ReportingDescriptor{}
	artifactLocations := []string{}

	tool := garif.NewDriver("protolint").
		WithInformationUri("https://github.com/yoheimuta/protolint")

	run := garif.NewRun(garif.NewTool(tool))

	for _, failure := range fs {
		_, ruleFound := rulesByID[failure.RuleID()]
		if !ruleFound {
			rule := garif.NewRule(
				failure.RuleID(),
			).
				WithHelpUri("https://github.com/yoheimuta/protolint")

			rulesByID[failure.RuleID()] = rule
			allRules = append(allRules, rule)
		}

		if !(contains(artifactLocations, failure.Pos().Filename)) {
			artifactLocations = append(artifactLocations, failure.Pos().Filename)
		}

		run.WithResult(
			failure.RuleID(),
			failure.Message(),
			failure.Pos().Filename,
			failure.Pos().Line,
			failure.Pos().Column,
		)
	}

	tool.WithRules(allRules...)
	run.WithArtifactsURIs(artifactLocations...)

	logFile := garif.NewLogFile([]*garif.Run{run}, garif.Version210)
	return logFile.PrettyWrite(w)
}
