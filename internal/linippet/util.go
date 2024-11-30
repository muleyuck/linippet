package linippet

import (
	"os"
	"path"
	"path/filepath"
)

const (
	ENV_NAME               = "LINIPPET_DATA"
	DEFAULT_LINIPPET_DIR   = ".linippet"
	DEFAULT_DATA_FILE_NAME = "linippet.json"
)

func getJsonPath() string {
	configPath, isExist := os.LookupEnv(ENV_NAME)
	if len(configPath) > 0 && isExist {
		// TODO: pathとして正しいか確認
		return configPath
	}
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, DEFAULT_LINIPPET_DIR, DEFAULT_DATA_FILE_NAME)
}

func checkJsonPath() (string, error) {
	dataPath := getJsonPath()
	baseDir := path.Dir(dataPath)
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		err = os.MkdirAll(baseDir, 0755)
		if err != nil {
			return "", err
		}
	}
	return dataPath, nil
}
