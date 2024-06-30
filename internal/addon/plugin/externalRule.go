package plugin

import (
	"path/filepath"

	"github.com/yoheimuta/protolint/internal/addon/plugin/shared"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/go-protoparser/v4/parser/meta"

	"github.com/yoheimuta/protolint/internal/addon/plugin/proto"
	"github.com/yoheimuta/protolint/linter/report"
	"github.com/yoheimuta/protolint/linter/rule"
)

// externalRule represents a customized rule that works as a plugin.
type externalRule struct {
	id       string
	purpose  string
	client   shared.RuleSet
	severity rule.Severity
}

func newExternalRule(
	id string,
	purpose string,
	client shared.RuleSet,
	severity rule.Severity,
) externalRule {
	return externalRule{
		id:       id,
		purpose:  purpose,
		client:   client,
		severity: severity,
	}
}

// ID returns the ID of this rule.
func (r externalRule) ID() string {
	return r.id
}

// Purpose returns the purpose of this rule.
func (r externalRule) Purpose() string {
	return r.purpose
}

// IsOfficial decides whether or not this rule belongs to the official guide.
func (r externalRule) IsOfficial() bool {
	return true
}

// Severity returns the severity of a rule (note, warning, error)
func (r externalRule) Severity() rule.Severity {
	return r.severity
}

// Apply applies the rule to the proto.
func (r externalRule) Apply(p *parser.Proto) ([]report.Failure, error) {
	relPath := p.Meta.Filename
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Apply(&proto.ApplyRequest{
		Id:   r.id,
		Path: absPath,
	})
	if err != nil {
		return nil, err
	}

	var fs []report.Failure
	for _, f := range resp.Failures {
		fs = append(fs, report.Failuref(meta.Position{
			Filename: relPath,
			Offset:   int(f.Pos.Offset),
			Line:     int(f.Pos.Line),
			Column:   int(f.Pos.Column),
		}, r.id, string(r.severity), f.Message))
	}
	return fs, nil
}
