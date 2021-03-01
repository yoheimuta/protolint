package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
) (*ExternalConfig, error) {
	newPath, err := getExternalConfigPath(filePath, dirPath)
	if err != nil {
		if len(filePath) == 0 && len(dirPath) == 0 {
			return nil, nil
		}
		return nil, err
	}

	data, err := ioutil.ReadFile(newPath)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}

	var config ExternalConfig
	if err := yaml.UnmarshalStrict(data, &config); err != nil {
		return nil, err
	}
	config.SourcePath = newPath

	return &config, nil
}

func getExternalConfigPath(
	filePath string,
	dirPath string,
) (string, error) {
	if 0 < len(filePath) {
		return filePath, nil
	}

	var checkedPaths []string
	for _, name := range []string{
		externalConfigFileName,
		externalConfigFileName2,
	} {
		for _, ext := range []string{
			externalConfigFileExtension,
			externalConfigFileExtension2,
		} {
			filePath := filepath.Join(dirPath, name+ext)
			checkedPaths = append(checkedPaths, filePath)
			if _, err := os.Stat(filePath); err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return "", err
			}
			return filePath, nil
		}
	}
	return "", fmt.Errorf("not found config file by searching `%s`", strings.Join(checkedPaths, ","))
}
