package config

import (
	yaml "gopkg.in/yaml.v2"
)

const (
	externalConfigFileName       = ".protolint"
	externalConfigFileName2      = "protolint"
	externalConfigFileExtension  = ".yaml"
	externalConfigFileExtension2 = ".yml"
)

type yamlConfigLoader struct {
	filePath string
}

func (y yamlConfigLoader) LoadExternalConfig() (*ExternalConfig, error) {
	data, err := loadFileContent(y.filePath)
	if err != nil {
		return nil, err
	}

	var config ExternalConfig

	if err := yaml.UnmarshalStrict(data, &config); err != nil {
		return nil, err
	}

	config.SourcePath = y.filePath

	return &config, nil
}
