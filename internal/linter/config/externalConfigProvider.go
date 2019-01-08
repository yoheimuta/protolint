package config

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

const (
	externalConfigFileName = "protolint.yaml"
)

// GetExternalConfig provides the externalConfig.
func GetExternalConfig(
	dirPath string,
) (ExternalConfig, error) {
	data, err := ioutil.ReadFile(filepath.Join(dirPath, externalConfigFileName))
	if err != nil {
		return ExternalConfig{}, err
	}
	if len(data) == 0 {
		return ExternalConfig{}, nil
	}

	var config ExternalConfig
	if err := yaml.UnmarshalStrict(data, &config); err != nil {
		return config, err
	}

	return config, nil
}
