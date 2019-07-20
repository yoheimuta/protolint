package config_test

import (
	"testing"

	"github.com/yoheimuta/protolint/internal/cmd/subcmds"
	"github.com/yoheimuta/protolint/internal/linter/config"
)

func TestExternalConfig_ShouldSkipRule(t *testing.T) {
	noDefaultExternalConfig := config.ExternalConfig{
		Lint: config.Lint{
			Ignores: []config.Ignore{
				{
					ID: "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
					Files: []string{
						"path/to/foo.proto",
						"path/to/bar.proto",
					},
				},
				{
					ID: "ENUM_NAMES_UPPER_CAMEL_CASE",
					Files: []string{
						"path/to/foo.proto",
					},
				},
			},
			Directories: config.Directories{
				Exclude: []string{
					"path/to/dir",
					"path/to/dir2",
				},
			},
			Rules: struct {
				NoDefault  bool     `yaml:"no_default"`
				AllDefault bool     `yaml:"all_default"`
				Add        []string `yaml:"add"`
				Remove     []string `yaml:"remove"`
			}{
				NoDefault: true,
				Add: []string{
					"FIELD_NAMES_LOWER_SNAKE_CASE",
					"MESSAGE_NAMES_UPPER_CAMEL_CASE",
				},
				Remove: []string{
					"RPC_NAMES_UPPER_CAMEL_CASE",
				},
			},
		},
	}

	defaultExternalConfig := config.ExternalConfig{
		Lint: config.Lint{
			Ignores: []config.Ignore{
				{
					ID: "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
					Files: []string{
						"path/to/foo.proto",
						"path/to/bar.proto",
					},
				},
				{
					ID: "ENUM_NAMES_UPPER_CAMEL_CASE",
					Files: []string{
						"path/to/foo.proto",
					},
				},
			},
			Rules: struct {
				NoDefault  bool     `yaml:"no_default"`
				AllDefault bool     `yaml:"all_default"`
				Add        []string `yaml:"add"`
				Remove     []string `yaml:"remove"`
			}{
				NoDefault: false,
				Add: []string{
					"FIELD_NAMES_LOWER_SNAKE_CASE",
					"MESSAGE_NAMES_UPPER_CAMEL_CASE",
				},
				Remove: []string{
					"RPC_NAMES_UPPER_CAMEL_CASE",
				},
			},
		},
	}

	for _, test := range []struct {
		name                string
		externalConfig      config.ExternalConfig
		inputRuleID         string
		inputDisplayPath    string
		inputDefaultRuleIDs []string
		wantSkipRule        bool
	}{
		{
			name:             "ignore ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
			externalConfig:   noDefaultExternalConfig,
			inputRuleID:      "ENUM_FIELD_NAMES_UPPER_SNAKE_CASE",
			inputDisplayPath: "path/to/foo.proto",
			wantSkipRule:     true,
		},
		{
			name:             "ignore ENUM_NAMES_UPPER_CAMEL_CASE",
			externalConfig:   noDefaultExternalConfig,
			inputRuleID:      "ENUM_NAMES_UPPER_CAMEL_CASE",
			inputDisplayPath: "path/to/foo.proto",
			wantSkipRule:     true,
		},
		{
			name:             "not ignore FIELD_NAMES_LOWER_SNAKE_CASE",
			externalConfig:   noDefaultExternalConfig,
			inputRuleID:      "FIELD_NAMES_LOWER_SNAKE_CASE",
			inputDisplayPath: "path/to/bar.proto",
		},
		{
			name:           "not skip Add rules",
			externalConfig: noDefaultExternalConfig,
			inputRuleID:    "FIELD_NAMES_LOWER_SNAKE_CASE",
		},
		{
			name:           "skip noAdd rules",
			externalConfig: noDefaultExternalConfig,
			inputRuleID:    "RPC_NAMES_UPPER_CAMEL_CASE",
			wantSkipRule:   true,
		},
		{
			name:           "skip Remove rule",
			externalConfig: defaultExternalConfig,
			inputRuleID:    "RPC_NAMES_UPPER_CAMEL_CASE",
			wantSkipRule:   true,
		},
		{
			name:           "not skip noRemove rule",
			externalConfig: defaultExternalConfig,
			inputRuleID:    "FIELD_NAMES_LOWER_SNAKE_CASE",
		},
		{
			name:           "not skip default rules",
			externalConfig: defaultExternalConfig,
			inputRuleID:    "HOGE_RULE",
			inputDefaultRuleIDs: []string{
				"HOGE_RULE",
			},
		},
		{
			name:                "not skip default one",
			externalConfig:      config.ExternalConfig{},
			inputRuleID:         subcmds.NewAllRules(config.RulesOption{}, false).Default().IDs()[0],
			inputDefaultRuleIDs: subcmds.NewAllRules(config.RulesOption{}, false).Default().IDs(),
		},
		{
			name:             "exclude the directory",
			externalConfig:   noDefaultExternalConfig,
			inputRuleID:      "FIELD_NAMES_LOWER_SNAKE_CASE",
			inputDisplayPath: "path/to/dir/bar.proto",
			wantSkipRule:     true,
		},
		{
			name:             "exclude the another directory",
			externalConfig:   noDefaultExternalConfig,
			inputRuleID:      "FIELD_NAMES_LOWER_SNAKE_CASE",
			inputDisplayPath: "path/to/dir2/bar.proto",
			wantSkipRule:     true,
		},
		{
			name:             "exclude the child directory",
			externalConfig:   noDefaultExternalConfig,
			inputRuleID:      "FIELD_NAMES_LOWER_SNAKE_CASE",
			inputDisplayPath: "path/to/dir2/child/bar.proto",
			wantSkipRule:     true,
		},
		{
			name:             "not exclude the another directory",
			externalConfig:   noDefaultExternalConfig,
			inputRuleID:      "FIELD_NAMES_LOWER_SNAKE_CASE",
			inputDisplayPath: "path/to/dir3/bar.proto",
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got := test.externalConfig.ShouldSkipRule(
				test.inputRuleID,
				test.inputDisplayPath,
				test.inputDefaultRuleIDs,
			)
			if got != test.wantSkipRule {
				t.Errorf("got %v, but want %v", got, test.wantSkipRule)
			}
		})
	}
}
