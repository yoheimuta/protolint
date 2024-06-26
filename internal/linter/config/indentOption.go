package config

import (
	"fmt"
	"strings"
)

// IndentOption represents the option for the INDENT rule.
type IndentOption struct {
	CustomizableSeverityOption
	Style string
	// Deprecated: not used
	Newline          string
	NotInsertNewline bool
}

// UnmarshalYAML implements yaml.v2 Unmarshaler interface.
func (i *IndentOption) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var option struct {
		CustomizableSeverityOption `yaml:",inline"`
		Style                      string `yaml:"style"`
		Newline                    string `yaml:"newline"`
		NotInsertNewline           bool   `yaml:"not_insert_newline"`
	}
	if err := unmarshal(&option); err != nil {
		return err
	}

	var style string
	switch option.Style {
	case "tab":
		style = "\t"
	case "4":
		style = strings.Repeat(" ", 4)
	case "2":
		style = strings.Repeat(" ", 2)
	case "":
		break
	default:
		return fmt.Errorf("%s is an invalid style option. valid option is tab, 4 or 2", option.Style)
	}
	i.Style = style

	switch option.Newline {
	case "\n", "\r", "\r\n", "":
		i.Newline = option.Newline
	default:
		return fmt.Errorf(`%s is an invalid newline option. valid option is \n, \r or \r\n`, option.Newline)
	}
	i.NotInsertNewline = option.NotInsertNewline
	i.CustomizableSeverityOption = option.CustomizableSeverityOption
	return nil
}

// UnmarshalTOML implements toml Unmarshaler interface.
func (i *IndentOption) UnmarshalTOML(data interface{}) error {
	optionsMap := map[string]interface{}{}
	for k, v := range data.(map[string]interface{}) {
		optionsMap[k] = v.(string)
	}

	if style, ok := optionsMap["style"]; ok {
		styleStr := style.(string)
		switch styleStr {
		case "\t":
			styleStr = "\t"
		case "tab":
			styleStr = "\t"
		case "4":
			styleStr = strings.Repeat(" ", 4)
		case "2":
			styleStr = strings.Repeat(" ", 2)
		case "":
			break
		default:
			return fmt.Errorf("%s is an invalid style option. valid option is \\t, tab, 4 or 2", style)
		}
		i.Style = styleStr
	}

	if newLine, ok := optionsMap["newline"]; ok {
		newLineStr := newLine.(string)
		switch newLineStr {
		case "\n", "\r", "\r\n", "":
			i.Newline = newLineStr
		default:
			return fmt.Errorf(`%s is an invalid newline option. valid option is \n, \r or \r\n`, newLine)
		}
	}

	if insertNoNewLine, ok := optionsMap["not_insert_newline"]; ok {
		i.NotInsertNewline = insertNoNewLine.(bool)
	}

	return nil
}
