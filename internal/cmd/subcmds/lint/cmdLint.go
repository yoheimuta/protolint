package lint

import (
	"fmt"
	"io"

	"github.com/yoheimuta/protolinter/internal/linter"
	"github.com/yoheimuta/protolinter/internal/linter/file"
	"github.com/yoheimuta/protolinter/internal/linter/report"
	"github.com/yoheimuta/protolinter/internal/osutil"
)

// CmdLint is a lint command.
type CmdLint struct {
	l          *linter.Linter
	stdout     io.Writer
	stderr     io.Writer
	protoFiles []file.ProtoFile
	config     CmdLintConfig
}

// NewCmdLint creates a new CmdLint.
func NewCmdLint(
	args []string,
	stdout io.Writer,
	stderr io.Writer,
) (*CmdLint, error) {
	protoSet, err := file.NewProtoSet(args)
	if err != nil {
		return nil, err
	}
	lintConfig := NewCmdLintConfig(protoSet.Config())

	return &CmdLint{
		l:          linter.NewLinter(),
		stdout:     stdout,
		stderr:     stderr,
		protoFiles: protoSet.ProtoFiles(),
		config:     lintConfig,
	}, nil
}

// Run lints to proto files.
func (c *CmdLint) Run() osutil.ExitCode {
	failures, err := c.run()
	if err != nil {
		_, _ = fmt.Fprintln(c.stderr, err)
		return osutil.ExitFailure
	}

	for _, failure := range failures {
		_, _ = fmt.Fprintln(c.stderr, failure)
	}
	if 0 < len(failures) {
		return osutil.ExitFailure
	}

	return osutil.ExitSuccess
}

func (c *CmdLint) run() ([]report.Failure, error) {
	var allFailures []report.Failure

	for _, f := range c.protoFiles {
		failures, err := c.runOneFile(f)
		if err != nil {
			return nil, err
		}
		allFailures = append(allFailures, failures...)
	}
	return allFailures, nil
}

func (c *CmdLint) runOneFile(
	f file.ProtoFile,
) ([]report.Failure, error) {
	proto, err := f.Parse()
	if err != nil {
		return nil, err
	}

	rs, err := c.config.GenRules(f)
	if err != nil {
		return nil, err
	}

	return c.l.Run(proto, rs)
}
