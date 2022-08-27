package autodisable

import "github.com/yoheimuta/go-protoparser/v4/parser"

// PlacementType is a selection of the placement strategies.
type PlacementType int

const (
	// Noop does nothing
	Noop PlacementType = iota
	// ThisThenNext puts inline comments.
	ThisThenNext
	// Next puts newline comments.
	Next
)

// NewPlacementStrategy creates a strategy object.
func NewPlacementStrategy(ptype PlacementType, filename, ruleID string) (PlacementStrategy, error) {
	if ptype == Noop {
		return &noopPlacementStrategy{}, nil
	}

	c, err := newCommentator(filename, ruleID)
	if err != nil {
		return nil, err
	}
	switch ptype {
	case ThisThenNext:
		return newThisThenNextPlacementStrategy(c), err
	case Next:
		return newNextPlacementStrategy(c), err
	default:
		return nil, nil
	}
}

// PlacementStrategy is an abstraction to put a comment.
type PlacementStrategy interface {
	Disable(offset int, comments []*parser.Comment, inline *parser.Comment)
	Finalize() error
}

type noopPlacementStrategy struct{}

func (p *noopPlacementStrategy) Disable(
	offset int,
	comments []*parser.Comment,
	_ *parser.Comment) {
}

func (p *noopPlacementStrategy) Finalize() error {
	return nil
}
