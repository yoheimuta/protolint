package list

import (
	"fmt"
	"io"

	"github.com/yoheimuta/protolint/internal/linter/config"

	"github.com/yoheimuta/protolint/internal/cmd/subcmds"
	"github.com/yoheimuta/protolint/internal/osutil"
	"github.com/yoheimuta/protolint/linter/rule"
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
		return osutil.ExitInternalFailure
	}
	return osutil.ExitSuccess
}

func (c *CmdList) run() error {
	rules, err := hasIDAndPurposes()
	if err != nil {
		return err
	}

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

func hasIDAndPurposes() ([]hasIDAndPurpose, error) {
	rs, err := subcmds.NewAllRules(config.RulesOption{}, false, false, nil)
	if err != nil {
		return nil, err
	}

	var rules []hasIDAndPurpose
	for _, r := range rs {
		rules = append(rules, r)
	}
	return rules, nil
}
