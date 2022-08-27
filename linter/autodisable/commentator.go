package autodisable

import (
	"log"

	"github.com/yoheimuta/go-protoparser/v4/parser"
	"github.com/yoheimuta/protolint/linter/disablerule"
	"github.com/yoheimuta/protolint/linter/fixer"
)

type commentator struct {
	fixing *fixer.BaseFixing
	ruleID string
}

func newCommentator(filename, ruleID string) (*commentator, error) {
	f, err := fixer.NewBaseFixing(filename)
	if err != nil {
		return nil, err
	}
	return &commentator{
		fixing: f,
		ruleID: ruleID,
	}, nil
}

func (c *commentator) insertNewline(offset int) {
	comment := disablerule.PrefixDisableNext + " " + c.ruleID

	space := ""
	pos := offset
loop:
	for i := offset; 0 < i; i-- {
		ch := c.fixing.Content()[i]
		switch ch {
		case ' ', '\t':
			space += " "
		case '\n', '\r':
			break loop
		default:
			pos = i
			space = ""
		}
	}
	c.insert("// "+comment+c.fixing.LineEnding()+space, pos)
}

func (c *commentator) tryMergeInline(inline *parser.Comment) bool {
	matches := disablerule.ReDisableThis.FindStringIndex(inline.Raw)
	log.Println(matches)
	if matches != nil {
		extracted := inline.Raw[matches[0]:matches[1]]
		log.Println(extracted)
		startPos := inline.Meta.Pos.Offset
		c.fixing.Replace(fixer.TextEdit{
			Pos:     startPos + matches[0],
			End:     startPos + matches[1] - 1,
			NewText: []byte(extracted + " " + c.ruleID),
		})
		return true
	}
	return false
}

func (c *commentator) insertInline(offset int) {
	comment := disablerule.PrefixDisableThis + " " + c.ruleID

	pos := offset
	content := c.fixing.Content()
loop:
	for i := offset; i < len(content); i++ {
		ch := content[i]
		switch ch {
		case ' ', '\t':
		case '\n', '\r':
			break loop
		default:
			pos = i
		}
	}
	c.insert(" // "+comment, pos+1)
}

func (c *commentator) finalize() error {
	return c.fixing.Finally()
}

func (c *commentator) insert(text string, pos int) {
	c.fixing.Replace(fixer.TextEdit{
		Pos:     pos,
		End:     pos - 1,
		NewText: []byte(text),
	})
}
