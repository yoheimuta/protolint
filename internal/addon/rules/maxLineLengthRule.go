package rules

import (
	"bufio"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/linter/disablerule"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
)

const (
	// Keep the line length to 80 characters.
	// See https://developers.google.com/protocol-buffers/docs/style#standard-file-formatting
	defaultMaxChars = 80

	defaultTabChars = 4
)

// MaxLineLengthRule enforces a maximum line length to increase code readability and maintainability.
// The length of a line is defined as the number of Unicode characters in the line.
type MaxLineLengthRule struct {
	RuleWithSeverity
	maxChars int
	tabChars int
}

// NewMaxLineLengthRule creates a new MaxLineLengthRule.
func NewMaxLineLengthRule(
	severity rule.Severity,
	maxChars int,
	tabChars int,
) MaxLineLengthRule {
	if maxChars == 0 {
		maxChars = defaultMaxChars
	}
	if tabChars == 0 {
		tabChars = defaultTabChars
	}
	return MaxLineLengthRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		maxChars:         maxChars,
		tabChars:         tabChars,
	}
}

// ID returns the ID of this rule.
func (r MaxLineLengthRule) ID() string {
	return "MAX_LINE_LENGTH"
}

// Purpose returns the purpose of this rule.
func (r MaxLineLengthRule) Purpose() string {
	return "Enforces a maximum line length."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r MaxLineLengthRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r MaxLineLengthRule) Apply(proto *parser.Proto) (
	failures []report.Failure,
	err error,
) {
	fileName := proto.Meta.Filename
	reader, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := reader.Close()
		if err != nil {
			return
		}
		if closeErr != nil {
			err = closeErr
		}
	}()

	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	disablerule.NewInterpreter(r.ID()).CallEachIfValid(
		lines,
		func(index int, line string) {
			line = strings.Replace(line, "\t", strings.Repeat(" ", r.tabChars), -1)
			lineCount := utf8.RuneCountInString(line)
			if r.maxChars < lineCount {
				failures = append(failures, report.Failuref(
					meta.Position{
						Filename: fileName,
						Line:     index + 1,
						Column:   1,
					},
					r.ID(),
					string(r.Severity()),
					"The line length is %d, but it must be shorter than %d",
					lineCount,
					r.maxChars,
				))
			}
		},
	)
	return failures, nil
}
