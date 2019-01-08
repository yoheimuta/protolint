package config

import (
	"io/ioutil"
	"os"
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
	filePath := filepath.Join(dirPath, externalConfigFileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ExternalConfig{}, nil
	}

	data, err := ioutil.ReadFile(filePath)
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
