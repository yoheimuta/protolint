package lint

import (
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/hashicorp/go-plugin"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/internal/linter/config"

	vfs "github.com/yoheimuta/protolint/internal/file"
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
	protoFiles []*file.ProtoFile
	config     CmdLintConfig
	output     io.Writer
}

// NewCmdLint creates a new CmdLint.
func NewCmdLint(
	flags Flags,
	stdout io.Writer,
	stderr io.Writer,
) (*CmdLint, error) {
	protoSet, err := file.NewProtoSet(flags.FilePaths, flags.StdinFilename)
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
	defer plugin.CleanupClients()

	failures, err := c.run()
	if err != nil {
		_, _ = fmt.Fprintln(c.stderr, err)
		return osutil.ExitInternalFailure
	}

	if c.config.IsModifyingMode() {
		for _, f := range c.protoFiles {
			if vfs.IsStdin(f.DisplayPath()) {
				fixedContent, err := vfs.ReadFile(f.DisplayPath())
				if err != nil {
					_, _ = fmt.Fprintln(c.stderr, "failed to read fixed stdin from VFS:", err)
					return osutil.ExitInternalFailure
				}

				_, _ = c.stdout.Write(fixedContent)

				return osutil.ExitSuccess
			}
		}
	}

	reportedFailures := failures
	if c.config.fixMode && len(failures) > 0 {
		remainingFailures, err := c.checkRemaining()
		if err != nil {
			_, _ = fmt.Fprintln(c.stderr, err)
			return osutil.ExitInternalFailure
		}
		reportedFailures = remainingFailures
	}

	err = c.config.reporters.ReportWithFallback(c.output, reportedFailures)
	if err != nil {
		_, _ = fmt.Fprintln(c.stderr, err)
		return osutil.ExitInternalFailure
	}

	if 0 < len(reportedFailures) {
		return osutil.ExitLintFailure
	}

	return osutil.ExitSuccess
}

func (c *CmdLint) run() ([]report.Failure, error) {
	var allFailures []report.Failure

	for i, f := range c.protoFiles {
		newF, failures, err := c.runOneFile(f, c.config)
		if err != nil {
			return nil, err
		}
		c.protoFiles[i] = newF
		allFailures = append(allFailures, failures...)
	}
	return allFailures, nil
}

func (c *CmdLint) checkRemaining() ([]report.Failure, error) {
	checkConfig := c.config.CheckConfig()
	var remainingFailures []report.Failure

	for _, f := range c.protoFiles {
		f.ResetCache()
		f.ResetData()

		_, failures, err := c.runOneFile(f, checkConfig)
		if err != nil {
			return nil, err
		}
		remainingFailures = append(remainingFailures, failures...)
	}
	return remainingFailures, nil
}

// ParseError represents the error returned through a parsing exception.
type ParseError struct {
	Message string
}

func (p ParseError) Error() string {
	return p.Message
}

func (c *CmdLint) runOneFile(
	f *file.ProtoFile,
	cfg CmdLintConfig,
) (*file.ProtoFile, []report.Failure, error) {
	// Gen rules first
	// If there is no rule, we can skip parse proto file
	rs, err := cfg.GenRules(f)
	if err != nil {
		return f, nil, err
	}
	if len(rs) == 0 {
		return f, []report.Failure{}, nil
	}

	failures, err := c.l.Run(func(p *parser.Proto) (*parser.Proto, error) {
		// Recreate a protoFile if the previous rule changed the filename.
		if p != nil && p.Meta.Filename != f.DisplayPath() {
			newFilename := p.Meta.Filename
			newBase := filepath.Base(newFilename)
			f = file.NewProtoFile(filepath.Join(filepath.Dir(f.Path()), newBase), newFilename)
		}

		if cfg.IsModifyingMode() {
			f.ResetCache()
			f.ResetData()
		}

		proto, err := f.Parse(cfg.verbose)
		if err != nil {
			if cfg.verbose {
				return nil, ParseError{Message: err.Error()}
			}
			return nil, ParseError{Message: fmt.Sprintf("%s. Use -v for more details", err)}
		}
		return proto, nil
	}, rs)
	if err != nil {
		return f, nil, err
	}
	return f, failures, nil
}
