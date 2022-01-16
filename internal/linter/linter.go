package linter

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
)

// Linter represents the protocol buffer linter with some rules.
type Linter struct{}

// NewLinter creates a new Linter.
func NewLinter() *Linter {
	return &Linter{}
}

// Run lints the protocol buffer.
func (l *Linter) Run(
	genProto func() (*parser.Proto, error),
	hasApplies []rule.HasApply,
) ([]report.Failure, error) {
	var fs []report.Failure
	for _, hasApply := range hasApplies {
		p, err := genProto()
		if err != nil {
			return nil, err
		}

		f, err := hasApply.Apply(p)
		if err != nil {
			return nil, err
		}
		fs = append(fs, f...)
	}
	return fs, nil
}
