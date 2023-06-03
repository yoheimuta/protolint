package config

// ServiceNamesEndWithOption represents the option for the SERVICE_NAMES_END_WITH rule.
type ServiceNamesEndWithOption struct {
	CustomizableSeverityOption
	Text string `yaml:"text"`
}
