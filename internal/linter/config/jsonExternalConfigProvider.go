package config

import (
	"encoding/json"
)

const packageJsonFileNameForJs = "package.json"
const packageJsonFileNameForJsExtension = ".json"

type jsonConfigLoader struct {
	filePath string
}

func (j jsonConfigLoader) LoadExternalConfig() (*ExternalConfig, error) {
	data, err := loadFileContent(j.filePath)
	if err != nil {
		return nil, err
	}

	var config ExternalConfig
	var jsonData jsonEmbeddedConfig
	// do not unmarshal strict. JS specific package.json will contain
	// other values as well.
	if jsonErr := json.Unmarshal(data, &jsonData); jsonErr != nil {
		return nil, jsonErr
	}

	readConfig := jsonData.toExternalConfig()
	if readConfig == nil {
		return nil, nil
	}
	config = *readConfig

	config.SourcePath = j.filePath

	return &config, nil
}

type jsonEmbeddedConfig struct {
	Protolint *Lint `json:"protolint"`
}

func (p jsonEmbeddedConfig) toExternalConfig() *ExternalConfig {
	if p.Protolint == nil {
		return nil
	}

	return &ExternalConfig{
		Lint: *p.Protolint,
	}
}
