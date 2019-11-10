package disablerule

import (
	"fmt"
	"regexp"
	"strings"
)

type commandType int

const (
	commandDisable commandType = iota
	commandEnable
	commandDisableNext
	commandDisableThis
)

var (
	reDisable     = regexp.MustCompile(`protolint:disable (.*)`)
	reEnable      = regexp.MustCompile(`protolint:enable (.*)`)
	reDisableNext = regexp.MustCompile(`protolint:disable:next (.*)`)
	reDisableThis = regexp.MustCompile(`protolint:disable:this (.*)`)
)

type command struct {
	ruleIDs []string
	t       commandType
}

func newCommand(
	comment string,
) (command, error) {
	subs := reDisable.FindStringSubmatch(comment)
	if len(subs) == 2 {
		ruleIDs := strings.Split(subs[1], " ")
		return command{
			ruleIDs: ruleIDs,
			t:       commandDisable,
		}, nil
	}

	subs = reEnable.FindStringSubmatch(comment)
	if len(subs) == 2 {
		ruleIDs := strings.Split(subs[1], " ")
		return command{
			ruleIDs: ruleIDs,
			t:       commandEnable,
		}, nil
	}

	subs = reDisableNext.FindStringSubmatch(comment)
	if len(subs) == 2 {
		ruleIDs := strings.Split(subs[1], " ")
		return command{
			ruleIDs: ruleIDs,
			t:       commandDisableNext,
		}, nil
	}

	subs = reDisableThis.FindStringSubmatch(comment)
	if len(subs) == 2 {
		ruleIDs := strings.Split(subs[1], " ")
		return command{
			ruleIDs: ruleIDs,
			t:       commandDisableThis,
		}, nil
	}

	return command{}, fmt.Errorf("invalid disabled comments")
}

func (c command) enabled(
	ruleID string,
) bool {
	return c.t == commandEnable && c.matchRuleID(ruleID)
}

func (c command) disabled(
	ruleID string,
) bool {
	return c.t == commandDisable && c.matchRuleID(ruleID)
}

func (c command) disabledNext(
	ruleID string,
) bool {
	return c.t == commandDisableNext && c.matchRuleID(ruleID)
}

func (c command) disabledThis(
	ruleID string,
) bool {
	return c.t == commandDisableThis && c.matchRuleID(ruleID)
}

func (c command) matchRuleID(
	ruleID string,
) bool {
	for _, id := range c.ruleIDs {
		if id == ruleID {
			return true
		}
	}
	return false
}
