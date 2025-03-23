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
		return filepath.Clean(configPath)
	}
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, DEFAULT_LINIPPET_DIR, DEFAULT_DATA_FILE_NAME)
}

func checkJsonPath() (dataPath string, err error) {
	dataPath = getJsonPath()
	// If data file not exists, create with initial value
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		baseDir := path.Dir(dataPath)
		err = os.MkdirAll(baseDir, 0755)
		if err != nil {
			return "", err
		}
		file, err := os.Create(dataPath)
		if err != nil {
			return "", err
		}
		defer func() {
			deferErr := file.Close()
			if deferErr != nil {
				err = deferErr
			}
		}()
		if _, err := file.Write([]byte("[]")); err != nil {
			return "", err
		}
	}
	return dataPath, nil
}
