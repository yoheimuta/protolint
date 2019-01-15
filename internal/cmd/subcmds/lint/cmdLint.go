package lint

import (
	"fmt"
	"io"

	"github.com/yoheimuta/protolint/internal/linter/config"

	"github.com/yoheimuta/protolint/internal/linter"
	"github.com/yoheimuta/protolint/internal/linter/file"
	"github.com/yoheimuta/protolint/internal/linter/report"
	"github.com/yoheimuta/protolint/internal/osutil"
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
	flags Flags,
	stdout io.Writer,
	stderr io.Writer,
) (*CmdLint, error) {
	protoSet, err := file.NewProtoSet(flags.Args())
	if err != nil {
		return nil, err
	}

	externalConfig, err := config.GetExternalConfig(flags.ConfigDirPath)
	if err != nil {
		return nil, err
	}
	lintConfig := NewCmdLintConfig(
		externalConfig,
		flags,
	)

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
	proto, err := f.Parse(c.config.verbose)
	if err != nil {
		if c.config.verbose {
			return nil, err
		}
		return nil, fmt.Errorf("%v. Use -v for more details.\n", err)
	}

	rs, err := c.config.GenRules(f)
	if err != nil {
		return nil, err
	}

	return c.l.Run(proto, rs)
}
