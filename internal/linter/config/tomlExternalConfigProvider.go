package config

import (
	"github.com/BurntSushi/toml"
)

const pyProjectTomlFileNameForPy = "pyproject.toml"
const pyProjectTomlFileNameForPyExtension = ".toml"

type tomlConfigLoader struct {
	filePath string
}

func (t tomlConfigLoader) LoadExternalConfig() (*ExternalConfig, error) {
	data, err := loadFileContent(t.filePath)
	if err != nil {
		return nil, err
	}

	var config ExternalConfig
	var tomlData tomlToolsEmbeddedConfig
	// do not unmarshal strict. JS specific package.json will contain
	// other values as well.
	if tomlErr := toml.Unmarshal(data, &tomlData); tomlErr != nil {
		return nil, tomlErr
	}

	readConfig := tomlData.toExternalConfig()
	if readConfig == nil {
		return nil, nil
	}
	config = *readConfig

	config.SourcePath = t.filePath

	return &config, nil
}

type tomlEmbeddedConfig struct {
	Protolint *Lint `toml:"protolint"`
}

type tomlToolsEmbeddedConfig struct {
	Tools *tomlEmbeddedConfig `toml:"tools"`
}

func (p tomlToolsEmbeddedConfig) toExternalConfig() *ExternalConfig {
	if p.Tools == nil {
		return nil
	}

	if p.Tools.Protolint == nil {
		return nil
	}

	return &ExternalConfig{
		Lint: *p.Tools.Protolint,
	}
}
