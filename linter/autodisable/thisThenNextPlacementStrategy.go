package autodisable

import "github.com/yoheimuta/go-protoparser/v4/parser"

type thisThenNextPlacementStrategy struct {
	c *commentator
}

func newThisThenNextPlacementStrategy(c *commentator) *thisThenNextPlacementStrategy {
	return &thisThenNextPlacementStrategy{
		c: c,
	}
}

func (p *thisThenNextPlacementStrategy) Disable(
	offset int,
	comments []*parser.Comment,
	inline *parser.Comment) {
	if inline == nil {
		p.c.insertInline(offset)
		return
	}
	if p.c.tryMergeInline(inline) {
		return
	}

	p.c.insertNewline(offset)
}

func (p *thisThenNextPlacementStrategy) Finalize() error {
	return p.c.finalize()
}
