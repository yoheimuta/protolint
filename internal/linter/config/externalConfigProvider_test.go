package config_test

import (
	"reflect"
	"testing"

	"github.com/yoheimuta/protolint/internal/setting_test"

	"github.com/yoheimuta/protolint/internal/linter/config"
)

func TestGetExternalConfig(t *testing.T) {
	for _, test := range []struct {
		name               string
		inputDirPath       string
		wantExternalConfig config.ExternalConfig
		wantExistErr       bool
	}{
		{
			name:         "invalid config file",
			inputDirPath: setting_test.TestDataPath("invalidconfig"),
			wantExistErr: true,
		},
		{
			name: "not found a config file",
		},
		{
			name:         "valid config file",
			inputDirPath: setting_test.TestDataPath("validconfig"),
			wantExternalConfig: config.ExternalConfig{
				Lint: struct {
					Ignores []struct {
						ID    string   `yaml:"id"`
						Files []string `yaml:"files"`
					}
					Rules struct {
						NoDefault bool     `yaml:"no_default"`
						Add       []string `yaml:"add"`
						Remove    []string `yaml:"remove"`
					}
					RuleOption config.RuleOption `yaml:"rule_option"`
				}{
					Ignores: []struct {
						ID    string   `yaml:"id"`
						Files []string `yaml:"files"`
					}{
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
						NoDefault bool     `yaml:"no_default"`
						Add       []string `yaml:"add"`
						Remove    []string `yaml:"remove"`
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
					RuleOption: config.RuleOption{
						MaxLineLength: config.MaxLineLengthOption{
							MaxChars: 80,
							TabChars: 2,
						},
					},
				},
			},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			got, err := config.GetExternalConfig(test.inputDirPath)
			if test.wantExistErr {
				if err == nil {
					t.Errorf("got err nil, but want err")
				}
				return
			}
			if err != nil {
				t.Errorf("got err %v, but want nil", err)
				return
			}

			if !reflect.DeepEqual(got, test.wantExternalConfig) {
				t.Errorf("got %v, but want %v", got, test.wantExternalConfig)
			}
		})
	}
}
