package rule

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
)

// Severity represents the severity of the rule.
// All failues will have this severity on export.
type Severity string

const (
	// SeverityNote represents a note only rule severity
	SeverityNote Severity = "note"
	// SeverityWarning represents a rule severity at a warning level
	SeverityWarning Severity = "warning"
	// SeverityError represents a rule severity at a warning level
	SeverityError Severity = "error"
)

// HasApply represents a rule which can be applied.
type HasApply interface {
	// Apply applies the rule to the proto.
	Apply(proto *parser.Proto) ([]report.Failure, error)
}

// HasID represents a rule with ID.
type HasID interface {
	// ID returns the ID of this rule. This should be all UPPER_SNAKE_CASE.
	ID() string
}

// HasPurpose represents a rule with Purpose.
type HasPurpose interface {
	// Purpose returns the purpose of this rule. This should be a human-readable string.
	Purpose() string
}

// HasIsOfficial represents a rule with IsOfficial.
type HasIsOfficial interface {
	// IsOfficial decides whether or not this rule belongs to the official guide.
	IsOfficial() bool
}

// HasSeverity represents a rule with a configurable severity
type HasSeverity interface {
	// Severity returns the selected severity of a rule
	Severity() Severity
}

// Rule represents a rule which a linter can apply.
type Rule interface {
	HasApply
	HasID
	HasPurpose
	HasIsOfficial
	HasSeverity
}
