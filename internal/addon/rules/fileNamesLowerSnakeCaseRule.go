package rules

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/yoheimuta/protolint/internal/stringsutil"

	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
	"github.com/yoheimuta/protolint/linter/strs"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// FileNamesLowerSnakeCaseRule verifies that all file names are lower_snake_case.proto.
// See https://developers.google.com/protocol-buffers/docs/style#file-structure.
type FileNamesLowerSnakeCaseRule struct {
	RuleWithSeverity
	excluded []string
	fixMode  bool
}

// NewFileNamesLowerSnakeCaseRule creates a new FileNamesLowerSnakeCaseRule.
func NewFileNamesLowerSnakeCaseRule(
	severity rule.Severity,
	excluded []string,
	fixMode bool,
) FileNamesLowerSnakeCaseRule {
	return FileNamesLowerSnakeCaseRule{
		RuleWithSeverity: RuleWithSeverity{severity: severity},
		excluded:         excluded,
		fixMode:          fixMode,
	}
}

// ID returns the ID of this rule.
func (r FileNamesLowerSnakeCaseRule) ID() string {
	return "FILE_NAMES_LOWER_SNAKE_CASE"
}

// Purpose returns the purpose of this rule.
func (r FileNamesLowerSnakeCaseRule) Purpose() string {
	return "Verifies that all file names are lower_snake_case.proto."
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r FileNamesLowerSnakeCaseRule) IsOfficial() bool {
	return true
}

// Apply applies the rule to the proto.
func (r FileNamesLowerSnakeCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &fileNamesLowerSnakeCaseVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(r.ID()),
		excluded:       r.excluded,
		fixMode:        r.fixMode,
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type fileNamesLowerSnakeCaseVisitor struct {
	*visitor.BaseAddVisitor
	excluded []string
	fixMode  bool
}

// OnStart checks the file.
func (v *fileNamesLowerSnakeCaseVisitor) OnStart(proto *parser.Proto) error {
	path := proto.Meta.Filename
	if stringsutil.ContainsStringInSlice(path, v.excluded) {
		return nil
	}

	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	if ext != ".proto" || !strs.IsLowerSnakeCase(base) {
		expected := strs.ToLowerSnakeCase(base) + ext
		v.AddFailurefWithProtoMeta(proto.Meta, "File name %q should be lower_snake_case.proto like %q.", filename, expected)

		if v.fixMode {
			dir := filepath.Dir(path)
			newPath := filepath.Join(dir, expected)
			if _, err := os.Stat(newPath); !os.IsNotExist(err) {
				v.AddFailurefWithProtoMeta(proto.Meta, "Failed to rename %q because %q already exists.", filename, expected)
				return nil
			}
			err := os.Rename(path, newPath)
			if err != nil {
				return err
			}

			// Notify the upstream this new filename by updating the proto.
			proto.Meta.Filename = newPath
		}
	}
	return nil
}
