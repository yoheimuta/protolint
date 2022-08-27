package autodisable

import "github.com/yoheimuta/go-protoparser/v4/parser"

type nextPlacementStrategy struct {
	c *commentator
}

func newNextPlacementStrategy(c *commentator) *nextPlacementStrategy {
	return &nextPlacementStrategy{
		c: c,
	}
}

func (p *nextPlacementStrategy) Disable(
	offset int,
	_ []*parser.Comment,
	_ *parser.Comment) {
	p.c.insertNewline(offset)
}

func (p *nextPlacementStrategy) Finalize() error {
	return p.c.finalize()
}
