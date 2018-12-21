package disablerule

import "github.com/yoheimuta/go-protoparser/parser"

// Interpreter represents an interpreter for disable rule comments.
type Interpreter struct {
	ruleID     string
	isDisabled bool
}

// NewInterpreter creates an Interpreter.
func NewInterpreter(
	ruleID string,
) *Interpreter {
	return &Interpreter{
		ruleID: ruleID,
	}
}

// Interpret interprets comments and returns a bool whether not apply the rule to a next element.
func (i *Interpreter) Interpret(
	comments []*parser.Comment,
) (isDisabledNext bool) {
	cmds := newCommands(comments)
	return i.interpret(cmds)
}

func (i *Interpreter) interpret(
	cmds commands,
) bool {
	id := i.ruleID
	if cmds.enabled(id) {
		i.isDisabled = false
		return false
	}
	if cmds.disabled(id) {
		i.isDisabled = true
		return true
	}
	if cmds.disabledNext(id) {
		return true
	}
	return i.isDisabled
}
