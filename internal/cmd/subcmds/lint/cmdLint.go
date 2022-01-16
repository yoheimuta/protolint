package lint

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/yoheimuta/go-protoparser/v4/parser"

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

	externalConfig, err := config.GetExternalConfig(flags.ConfigPath, flags.ConfigDirPath)
	if err != nil {
		return nil, err
	}
	if flags.Verbose {
		if externalConfig != nil {
			log.Printf("[INFO] protolint loads a config file at %s\n", externalConfig.SourcePath)
		} else {
			log.Println("[INFO] protolint doesn't load a config file")
		}
	}
	if externalConfig == nil {
		externalConfig = &(config.ExternalConfig{})
	}
	lintConfig := NewCmdLintConfig(
		*externalConfig,
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
		return osutil.ExitInternalFailure
	}

	err = c.config.reporter.Report(c.output, failures)
	if err != nil {
		_, _ = fmt.Fprintln(c.stderr, err)
		return osutil.ExitInternalFailure
	}

	if 0 < len(failures) {
		return osutil.ExitLintFailure
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

// ParseError represents the error returned through a parsing exception.
type ParseError struct {
	Message string
}

func (p ParseError) Error() string {
	return p.Message
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

	return c.l.Run(func() (*parser.Proto, error) {
		proto, err := f.Parse(c.config.verbose)
		if err != nil {
			if c.config.verbose {
				return nil, ParseError{Message: err.Error()}
			}
			return nil, ParseError{Message: fmt.Sprintf("%s. Use -v for more details", err)}
		}
		return proto, nil
	}, rs)
}
