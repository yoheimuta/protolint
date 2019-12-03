package rule

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"

	"github.com/yoheimuta/protolint/linter/report"
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

// Rule represents a rule which a linter can apply.
type Rule interface {
	HasApply
	HasID
	HasPurpose
	HasIsOfficial
}
