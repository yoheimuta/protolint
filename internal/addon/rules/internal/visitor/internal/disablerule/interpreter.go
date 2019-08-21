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

// Interpret interprets comments and returns a bool whether not apply the rule to a next or this element.
func (i *Interpreter) Interpret(
	comments []*parser.Comment,
	inlines ...*parser.Comment,
) (isDisabled bool) {
	cmds := newCommands(comments)
	inlineCmds := newCommands(inlines)
	allCmds := append(append([]command{}, cmds...), inlineCmds...)
	return i.interpret(allCmds) ||
		i.interpretNext(cmds) ||
		i.interpretThis(inlineCmds) ||
		i.isDisabled
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
	return false
}

func (i *Interpreter) interpretNext(
	cmds commands,
) bool {
	id := i.ruleID
	if cmds.disabledNext(id) {
		return true
	}
	return false
}

func (i *Interpreter) interpretThis(
	cmds commands,
) bool {
	id := i.ruleID
	if cmds.disabledThis(id) {
		return true
	}
	return false
}
