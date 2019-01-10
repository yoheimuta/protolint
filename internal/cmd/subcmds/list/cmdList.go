package list

import (
	"io"

	"github.com/yoheimuta/protolint/internal/linter/config"

	"fmt"

	"github.com/yoheimuta/protolint/internal/cmd/subcmds"
	"github.com/yoheimuta/protolint/internal/linter/rule"
	"github.com/yoheimuta/protolint/internal/osutil"
)

// CmdList is a rule list command.
type CmdList struct {
	stdout io.Writer
	stderr io.Writer
}

// NewCmdList creates a new CmdList.
func NewCmdList(
	stdout io.Writer,
	stderr io.Writer,
) *CmdList {
	return &CmdList{
		stdout: stdout,
		stderr: stderr,
	}
}

// Run lists each rule description.
func (c *CmdList) Run() osutil.ExitCode {
	err := c.run()
	if err != nil {
		return osutil.ExitFailure
	}
	return osutil.ExitSuccess
}

func (c *CmdList) run() error {
	rules := hasIDAndPurposes()
	for _, r := range rules {
		_, err := fmt.Fprintf(
			c.stdout,
			"%s: %s\n",
			r.ID(),
			r.Purpose(),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

type hasIDAndPurpose interface {
	rule.HasID
	rule.HasPurpose
}

func hasIDAndPurposes() []hasIDAndPurpose {
	var rules []hasIDAndPurpose
	for _, r := range subcmds.NewAllRules(config.RulesOption{}, false) {
		rules = append(rules, r)
	}
	return rules
}
