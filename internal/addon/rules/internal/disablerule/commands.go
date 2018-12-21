package disablerule

import (
	"github.com/yoheimuta/go-protoparser/parser"
)

type commands []command

func newCommands(
	comments []*parser.Comment,
) commands {
	var cmds []command
	for _, comment := range comments {
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
