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

// comment prefix
const (
	PrefixDisable     = `protolint:disable`
	PrefixEnable      = `protolint:enable`
	PrefixDisableNext = `protolint:disable:next`
	PrefixDisableThis = `protolint:disable:this`
)

// comment prefix regexp
var (
	ReDisable     = regexp.MustCompile(PrefixDisable + ` (.*)`)
	ReEnable      = regexp.MustCompile(PrefixEnable + ` (.*)`)
	ReDisableNext = regexp.MustCompile(PrefixDisableNext + ` (.*)`)
	ReDisableThis = regexp.MustCompile(PrefixDisableThis + ` (.*)`)
)

type command struct {
	ruleIDs []string
	t       commandType
}

func newCommand(
	comment string,
) (command, error) {
	subs := ReDisable.FindStringSubmatch(comment)
	if len(subs) == 2 {
		ruleIDs := strings.Fields(strings.TrimSpace(subs[1]))
		return command{
			ruleIDs: ruleIDs,
			t:       commandDisable,
		}, nil
	}

	subs = ReEnable.FindStringSubmatch(comment)
	if len(subs) == 2 {
		ruleIDs := strings.Fields(strings.TrimSpace(subs[1]))
		return command{
			ruleIDs: ruleIDs,
			t:       commandEnable,
		}, nil
	}

	subs = ReDisableNext.FindStringSubmatch(comment)
	if len(subs) == 2 {
		ruleIDs := strings.Fields(strings.TrimSpace(subs[1]))
		return command{
			ruleIDs: ruleIDs,
			t:       commandDisableNext,
		}, nil
	}

	subs = ReDisableThis.FindStringSubmatch(comment)
	if len(subs) == 2 {
		ruleIDs := strings.Fields(strings.TrimSpace(subs[1]))
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
