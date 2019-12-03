package disablerule

import (
	"github.com/yoheimuta/go-protoparser/v4/parser"
)

type commands []command

func newCommands(
	comments []*parser.Comment,
) commands {
	var cmds []command
	for _, comment := range comments {
		if comment == nil {
			continue
		}
		cmd, err := newCommand(comment.Raw)
		if err == nil {
			cmds = append(cmds, cmd)
		}
	}
	return cmds
}

func (cs commands) enabled(
	ruleID string,
) bool {
	for _, cmd := range cs {
		if cmd.enabled(ruleID) {
			return true
		}
	}
	return false
}

func (cs commands) disabled(
	ruleID string,
) bool {
	for _, cmd := range cs {
		if cmd.disabled(ruleID) {
			return true
		}
	}
	return false
}

func (cs commands) disabledNext(
	ruleID string,
) bool {
	for _, cmd := range cs {
		if cmd.disabledNext(ruleID) {
			return true
		}
	}
	return false
}

func (cs commands) disabledThis(
	ruleID string,
) bool {
	for _, cmd := range cs {
		if cmd.disabledThis(ruleID) {
			return true
		}
	}
	return false
}
