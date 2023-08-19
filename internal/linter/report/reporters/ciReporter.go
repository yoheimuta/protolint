package reporters

import (
	"bytes"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/yoheimuta/protolint/linter/report"
)

type CiPipelineLogTemplate string

const (
	// problemMatcher provides a generic issue line that can be parsed in a problem matcher or jenkins pipeline
	problemMatcher CiPipelineLogTemplate = "Protolint {{ .Rule }} ({{ .Severity }}): {{ .File }}[{{ .Line }},{{ .Column }}]: {{ .Message }}"
	// azureDevOps provides a log message according to https://learn.microsoft.com/en-us/azure/devops/pipelines/scripts/logging-commands?view=azure-devops&tabs=bash#task-commands
	azureDevOps CiPipelineLogTemplate = "{{ if ne \"info\" .Severity }}##vso[task.logissue type={{ .Severity }};sourcepath={{ .File }};linenumber={{ .Line }};columnnumber={{ .Column }};code={{ .Rule }};]{{ .Message }}{{end}}"
	// gitlabCiCd provides an issue template where the severity is written in upper case. This is matched by the pipeline
	gitlabCiCd CiPipelineLogTemplate = "{{ .Severity | ToUpper }}: {{ .Rule }}  {{ .File }}({{ .Line }},{{ .Column }}) : {{ .Message }}"
	// githubActions provides an issue template according to https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#setting-a-notice-message
	githubActions CiPipelineLogTemplate = "::{{ if ne \"info\" .Severity }}{{ .Severity }}{{ else }}notice{{ end }} file={{ .File }},line={{ .Line }},col={{ .Column }},title={{ .Rule }}::{{ .Message }}"
	// empty provides default value for invalid returns
	empty CiPipelineLogTemplate = ""
	// env provides a marker for processing CI Templates from environment
	env CiPipelineLogTemplate = "[ENV]"
)

type CiReporter struct {
	pattern CiPipelineLogTemplate
}

func NewCiReporterForAzureDevOps() CiReporter {
	return CiReporter{pattern: azureDevOps}
}

func NewCiReporterForGitlab() CiReporter {
	return CiReporter{pattern: gitlabCiCd}
}

func NewCiReporterForGithubActions() CiReporter {
	return CiReporter{pattern: githubActions}
}

func NewCiReporterWithGenericFormat() CiReporter {
	return CiReporter{pattern: problemMatcher}
}

func NewCiReporterFromEnv() CiReporter {
	return CiReporter{pattern: env}
}

type ciReportedFailure struct {
	Severity string
	File     string
	Line     int
	Column   int
	Rule     string
	Message  string
}

func (c CiReporter) Report(w io.Writer, fs []report.Failure) error {
	template, err := c.getTemplate()
	if err != nil {
		return err
	}

	for _, failure := range fs {
		reportedFailure := ciReportedFailure{
			Severity: getSeverity(failure.Severity()),
			File:     failure.Pos().Filename,
			Line:     failure.Pos().Line,
			Column:   failure.Pos().Column,
			Rule:     failure.RuleID(),
			Message:  strings.Trim(strconv.Quote(failure.Message()), `"`), // Ensure message is on a single line without quotes
		}

		var buffer bytes.Buffer
		err = template.Execute(&buffer, reportedFailure)
		if err != nil {
			return err
		}

		written, err := w.Write(buffer.Bytes())
		if err != nil {
			return err
		}

		if written > 0 {
			w.Write([]byte("\n"))
		}
	}

	return nil
}

func getSeverity(s string) string {
	if s == "note" {
		return "info"
	}

	return s
}

func (c CiReporter) getTemplateString() CiPipelineLogTemplate {
	if c.pattern == env {
		template, err := getPatternFromEnv()
		if err != nil {
			log.Printf("[ERROR] Failed to process template from Environment: %s\n", err.Error())
			return problemMatcher
		}

		if template == empty {
			return problemMatcher
		}

		return template
	}
	if c.pattern != empty {
		return c.pattern
	}
	return problemMatcher
}

func (c CiReporter) getTemplate() (*template.Template, error) {
	toParse := c.getTemplateString()

	toUpper := template.FuncMap{"ToUpper": strings.ToUpper}

	template := template.New("Failure").Funcs(toUpper)
	evaluate, err := template.Parse(string(toParse))
	if err != nil {
		return nil, err
	}
	return evaluate, nil
}

func getPatternFromEnv() (CiPipelineLogTemplate, error) {
	templateString := os.Getenv("PROTOLINT_CIREPORTER_TEMPLATE_STRING")
	if templateString != "" {
		return CiPipelineLogTemplate(templateString), nil
	}

	templateFile := os.Getenv("PROTOLINT_CIREPORTER_TEMPLATE_FILE")
	if templateFile != "" {
		content, err := os.ReadFile(templateFile)

		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("[ERROR] Failed to open file %s from 'PROTOLINT_CIREPORTER_TEMPLATE_FILE'. File does not exist.\n", templateFile)
				log.Println("[WARN] Starting output with default processor.")
				return empty, nil
			}
			if os.IsPermission(err) {
				log.Printf("[ERROR] Failed to open file %s from 'PROTOLINT_CIREPORTER_TEMPLATE_FILE'. Insufficient permissions.\n", templateFile)
				log.Println("[WARN] Starting output with default processor.")
				return empty, nil
			}

			return empty, err
		}

		content_string := string(content)

		return CiPipelineLogTemplate(content_string), nil
	}

	return empty, nil
}
