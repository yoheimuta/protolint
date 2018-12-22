package rules

import (
	"github.com/yoheimuta/go-protoparser/parser"
	"github.com/yoheimuta/protolint/internal/addon/rules/internal/visitor"
	"github.com/yoheimuta/protolint/internal/linter/report"
	"github.com/yoheimuta/protolint/internal/strs"
)

// ServiceNamesUpperCamelCaseRule verifies that all service names are CamelCase (with an initial capital).
// See https://developers.google.com/protocol-buffers/docs/style#services.
type ServiceNamesUpperCamelCaseRule struct{}

// NewServiceNamesUpperCamelCaseRule creates a new ServiceNamesUpperCamelCaseRule.
func NewServiceNamesUpperCamelCaseRule() ServiceNamesUpperCamelCaseRule {
	return ServiceNamesUpperCamelCaseRule{}
}

// ID returns the ID of this rule.
func (r ServiceNamesUpperCamelCaseRule) ID() string {
	return "SERVICE_NAMES_UPPER_CAMEL_CASE"
}

// Purpose returns the purpose of this rule.
func (r ServiceNamesUpperCamelCaseRule) Purpose() string {
	return "Verifies that all service names are CamelCase (with an initial capital)."
}

// Apply applies the rule to the proto.
func (r ServiceNamesUpperCamelCaseRule) Apply(proto *parser.Proto) ([]report.Failure, error) {
	v := &serviceNamesUpperCamelCaseVisitor{
		BaseAddVisitor: visitor.NewBaseAddVisitor(),
	}
	return visitor.RunVisitor(v, proto, r.ID())
}

type serviceNamesUpperCamelCaseVisitor struct {
	*visitor.BaseAddVisitor
}

// VisitService checks the service.
func (v *serviceNamesUpperCamelCaseVisitor) VisitService(service *parser.Service) bool {
	if !strs.IsUpperCamelCase(service.ServiceName) {
		v.AddFailuref(service.Meta.Pos, "Service name %q must be UpperCamelCase", service.ServiceName)
	}
	return false
}
