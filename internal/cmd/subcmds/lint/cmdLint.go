package lint

import (
	"fmt"
	"io"
	"os"

	"github.com/yoheimuta/protolint/internal/linter/config"

	"github.com/yoheimuta/protolint/internal/linter"
	"github.com/yoheimuta/protolint/internal/linter/file"
	"github.com/yoheimuta/protolint/internal/osutil"
	"github.com/yoheimuta/protolint/linter/report"
)

// CmdLint is a lint command.
type CmdLint struct {
	l          *linter.Linter
	stdout     io.Writer
	stderr     io.Writer
	protoFiles []file.ProtoFile
	config     CmdLintConfig
	output     io.Writer
}

// NewCmdLint creates a new CmdLint.
func NewCmdLint(
	flags Flags,
	stdout io.Writer,
	stderr io.Writer,
) (*CmdLint, error) {
	protoSet, err := file.NewProtoSet(flags.FilePaths)
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

	output := stderr
	if 0 < len(flags.OutputFilePath) {
		output, err = os.OpenFile(flags.OutputFilePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
	}

	return &CmdLint{
		l:          linter.NewLinter(),
		stdout:     stdout,
		stderr:     stderr,
		protoFiles: protoSet.ProtoFiles(),
		config:     lintConfig,
		output:     output,
	}, nil
}

// Run lints to proto files.
func (c *CmdLint) Run() osutil.ExitCode {
	failures, err := c.run()
	if err != nil {
		_, _ = fmt.Fprintln(c.stderr, err)
		return osutil.ExitFailure
	}

	err = c.config.reporter.Report(c.output, failures)
	if err != nil {
		_, _ = fmt.Fprintln(c.stderr, err)
		return osutil.ExitFailure
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
	// Gen rules first
	// If there is no rule, we can skip parse proto file
	rs, err := c.config.GenRules(f)
	if err != nil {
		return nil, err
	}
	if len(rs) == 0 {
		return []report.Failure{}, nil
	}

	proto, err := f.Parse(c.config.verbose)
	if err != nil {
		if c.config.verbose {
			return nil, err
		}
		return nil, fmt.Errorf("%s. Use -v for more details", err)
	}

	return c.l.Run(proto, rs)
}
