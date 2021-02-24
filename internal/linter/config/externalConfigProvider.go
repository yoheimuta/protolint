package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

const (
	externalConfigFileName       = ".protolint"
	externalConfigFileName2      = "protolint"
	externalConfigFileExtension  = ".yaml"
	externalConfigFileExtension2 = ".yml"
)

// GetExternalConfig provides the externalConfig.
func GetExternalConfig(
	filePath string,
	dirPath string,
) (ExternalConfig, error) {
	filePath, err := getExternalConfigPath(filePath, dirPath)
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
	filePath string,
	dirPath string,
) (string, error) {
	if 0 < len(filePath) {
		return filePath, nil
	}
	for _, name := range []string{
		externalConfigFileName,
		externalConfigFileName2,
	} {
		for _, ext := range []string{
			externalConfigFileExtension,
			externalConfigFileExtension2,
		} {
			filePath := filepath.Join(dirPath, name+ext)
			if _, err := os.Stat(filePath); err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return "", err
			}
			return filePath, nil
		}
	}
	return "", nil
}
