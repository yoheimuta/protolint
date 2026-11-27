package config

// RPCRequestResponseNamesEndWithOption represents the option for the RPC_REQUEST_RESPONSE_NAMES_END_WITH rule.
type RPCRequestResponseNamesEndWithOption struct {
	CustomizableSeverityOption `yaml:",inline"`
	RequestSuffix              string   `yaml:"request_suffix" json:"request_suffix" toml:"request_suffix"`
	ResponseSuffix             string   `yaml:"response_suffix" json:"response_suffix" toml:"response_suffix"`
	Excludes                   []string `yaml:"excludes" json:"excludes" toml:"excludes"`
}
