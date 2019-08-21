package rules

import (
	"path/filepath"
	"strings"

	"github.com/yoheimuta/protolint/internal/stringsutil"

	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/strs"
	"github.com/yoheimuta/protolint/linter/visitor"
)

// FileNamesLowerSnakeCaseRule verifies that all file names are lower_snake_case.proto.
// See https://developers.google.com/protocol-buffers/docs/style#file-structure.
type FileNamesLowerSnakeCaseRule struct {
	excluded []string
}

// NewFileNamesLowerSnakeCaseRule creates a new FileNamesLowerSnakeCaseRule.
func NewFileNamesLowerSnakeCaseRule(
	excluded []string,
) FileNamesLowerSnakeCaseRule {
	return FileNamesLowerSnakeCaseRule{
		excluded: excluded,
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
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type fileNamesLowerSnakeCaseVisitor struct {
	*visitor.BaseAddVisitor
	excluded []string
}

// OnStart checks the file.
func (v *fileNamesLowerSnakeCaseVisitor) OnStart(proto *parser.Proto) error {
	path := proto.Meta.Filename
	if stringsutil.ContainsStringInSlice(path, v.excluded) {
		return nil
	}

	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	if ext != ".proto" || !strs.IsLowerSnakeCase(strings.TrimSuffix(filename, ext)) {
		v.AddFailurefWithProtoMeta(proto.Meta, "File name should be lower_snake_case.proto.")
	}
	return nil
}
