package plugin

import (
	"fmt"

	"github.com/yoheimuta/protolint/internal/addon/plugin/proto"
	"github.com/yoheimuta/protolint/internal/linter/file"
	"github.com/yoheimuta/protolint/internal/linter/rule"
)

type ruleSet struct {
	rawRules []rule.Rule

	rules   map[string]rule.Rule
	verbose bool
}

func newRuleSet(rules []rule.Rule) *ruleSet {
	return &ruleSet{
		rawRules: rules,
	}
}

func (c *ruleSet) initialize(req *proto.ListRulesRequest) {
	c.verbose = req.Verbose

	ruleMap := make(map[string]rule.Rule)
	for _, r := range c.rawRules {
		if f, ok := r.(RuleGen); ok {
			r = f(
				req.Verbose,
				req.FixMode,
			)
		}
		ruleMap[r.ID()] = r
	}
	c.rules = ruleMap
}

func (c *ruleSet) ListRules(req *proto.ListRulesRequest) (*proto.ListRulesResponse, error) {
	c.initialize(req)

	var meta []*proto.ListRulesResponse_Rule
	for _, r := range c.rules {
		meta = append(meta, &proto.ListRulesResponse_Rule{
			Id:      r.ID(),
			Purpose: r.Purpose(),
		})
	}
	return &proto.ListRulesResponse{
		Rules: meta,
	}, nil
}

func (c *ruleSet) Apply(req *proto.ApplyRequest) (*proto.ApplyResponse, error) {
	r, ok := c.rules[req.Id]
	if !ok {
		return nil, fmt.Errorf("not found rule=%s", req.Id)
	}

	absPath := req.Path
	protoFile := file.NewProtoFile(absPath, absPath)
	p, err := protoFile.Parse(c.verbose)
	if err != nil {
		return nil, err
	}

	fs, err := r.Apply(p)
	if err != nil {
		return nil, err
	}
	var fsp []*proto.ApplyResponse_Failure
	for _, f := range fs {
		fsp = append(fsp, &proto.ApplyResponse_Failure{
			Message: f.Message(),
			Pos: &proto.ApplyResponse_Position{
				Offset: int32(f.Pos().Offset),
				Line:   int32(f.Pos().Line),
				Column: int32(f.Pos().Column),
			},
		})
	}
	return &proto.ApplyResponse{
		Failures: fsp,
	}, nil
}
