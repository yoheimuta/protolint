package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/yoheimuta/protolint/internal/stringsutil"

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
		return nil, fmt.Errorf("read %s, but the content is empty", newPath)
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

	dirPaths := []string{dirPath}
	if len(dirPath) == 0 {
		absPath, err := os.Getwd()
		if err != nil {
			return "", err
		}

		absPath = filepath.Dir(absPath)
		for !stringsutil.ContainsStringInSlice(absPath, dirPaths) {
			dirPaths = append(dirPaths, absPath)
			absPath = filepath.Dir(absPath)
		}
	}

	var checkedPaths []string
	for _, dir := range dirPaths {
		for _, name := range []string{
			externalConfigFileName,
			externalConfigFileName2,
		} {
			for _, ext := range []string{
				externalConfigFileExtension,
				externalConfigFileExtension2,
			} {
				filePath := filepath.Join(dir, name+ext)
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
	}
	return "", fmt.Errorf("not found config file by searching `%s`", strings.Join(checkedPaths, ","))
}
