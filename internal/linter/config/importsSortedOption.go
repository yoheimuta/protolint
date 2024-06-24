package config

import (
	"fmt"
)

// ImportsSortedOption represents the option for the IMPORTS_SORTED rule.
type ImportsSortedOption struct {
	CustomizableSeverityOption
	// Deprecated: not used
	Newline string
}

// UnmarshalYAML implements yaml.v2 Unmarshaler interface.
func (i *ImportsSortedOption) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var option struct {
		CustomizableSeverityOption `yaml:",inline"`
		Newline                    string `yaml:"newline"`
	}
	if err := unmarshal(&option); err != nil {
		return err
	}

	switch option.Newline {
	case "\n", "\r", "\r\n", "":
		i.Newline = option.Newline
	default:
		return fmt.Errorf(`%s is an invalid newline option. valid option is \n, \r or \r\n`, option.Newline)
	}
	i.CustomizableSeverityOption = option.CustomizableSeverityOption
	return nil
}

// UnmarshalTOML implements toml Unmarshaler interface.
func (i *ImportsSortedOption) UnmarshalTOML(data interface{}) error {
	optionsMap := map[string]interface{}{}
	for k, v := range data.(map[string]interface{}) {
		optionsMap[k] = v.(string)
	}

	if newline, ok := optionsMap["newline"]; ok {
		switch newline.(string) {
		case "\n", "\r", "\r\n", "":
			i.Newline = newline.(string)
		default:
			return fmt.Errorf(`%s is an invalid newline option. valid option is \n, \r or \r\n`, newline)
		}
	}
	return nil
}
