package linter

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/internal/linter/report"
	"github.com/yoheimuta/protolint/internal/linter/rule"
)

// Linter represents the protocol buffer linter with some rules.
type Linter struct{}

// NewLinter creates a new Linter.
func NewLinter() *Linter {
	return &Linter{}
}

// Run lints the protocol buffer.
func (l *Linter) Run(
	proto *parser.Proto,
	hasApplies []rule.HasApply,
) ([]report.Failure, error) {
	var fs []report.Failure
	for _, hasApply := range hasApplies {
		f, err := hasApply.Apply(proto)
		if err != nil {
			return nil, err
		}
		fs = append(fs, f...)
	}
	return fs, nil
}
