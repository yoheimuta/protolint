package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

const (
	externalConfigFileName  = ".protolint.yaml"
	externalConfigFileName2 = "protolint.yaml"
)

// GetExternalConfig provides the externalConfig.
func GetExternalConfig(
	dirPath string,
) (ExternalConfig, error) {
	filePath, err := getExternalConfigPath(dirPath)
	if err != nil {
		return ExternalConfig{}, err
	}
	if len(filePath) == 0 {
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

func getExternalConfigPath(
	dirPath string,
) (string, error) {
	for _, name := range []string{
		externalConfigFileName,
		externalConfigFileName2,
	} {
		filePath := filepath.Join(dirPath, name)
		if _, err := os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}
		return filePath, nil
	}
	return "", nil
}
