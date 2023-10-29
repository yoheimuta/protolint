package disablerule

import "github.com/yoheimuta/go-protoparser/v4/parser"

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

// CallEachIfValid calls a given function each time the line is not disabled.
func (i *Interpreter) CallEachIfValid(
	lines []string,
	f func(index int, line string),
) {
	shouldSkip := false

	for index, line := range lines {
		cmd, err := newCommand(line)
		if err != nil {
			if !i.isDisabled && !shouldSkip {
				f(index, line)
			}
			if shouldSkip {
				shouldSkip = false
			}
			continue
		}

		if cmd.enabled(i.ruleID) {
			i.isDisabled = false
			f(index, line)
			continue
		}

		if cmd.disabled(i.ruleID) {
			i.isDisabled = true
			continue
		}

		if cmd.disabledThis(i.ruleID) {
			continue
		}

		if cmd.disabledNext(i.ruleID) {
			shouldSkip = true
			f(index, line)
			continue
		}

		if shouldSkip {
			shouldSkip = false
			continue
		}
		if i.isDisabled {
			continue
		}

		f(index, line)
	}
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
