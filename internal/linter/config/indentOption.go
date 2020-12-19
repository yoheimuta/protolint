package config

import (
	"fmt"
	"strings"
)

// IndentOption represents the option for the INDENT rule.
type IndentOption struct {
	Style            string
	Newline          string
	NotInsertNewline bool
}

// UnmarshalYAML implements yaml.v2 Unmarshaler interface.
func (i *IndentOption) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var option struct {
		Style            string `yaml:"style"`
		Newline          string `yaml:"newline"`
		NotInsertNewline bool   `yaml:"not_insert_newline"`
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
	return nil
}
