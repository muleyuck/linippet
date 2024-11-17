package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

const (
	ENV_NAME               = "LINIPPET_DATA"
	DEFAULT_LINIPPET_DIR   = ".linippet"
	DEFAULT_DATA_FILE_NAME = "linippet.json"
)

func getDataPath() string {
	configPath, isExist := os.LookupEnv(ENV_NAME)
	if len(configPath) > 0 && isExist {
		// TODO: pathとして正しいか確認
		return configPath
	}
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, DEFAULT_LINIPPET_DIR, DEFAULT_DATA_FILE_NAME)
}

func CheckDataPath() (string, error) {
	dataPath := getDataPath()
	baseDir := path.Dir(dataPath)
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		err = os.MkdirAll(baseDir, 0755)
		if err != nil {
			return "", err
		}
	}
	return dataPath, nil
}

type LinippetData struct {
	Snippet string `json:"snippet"`
}

func ReadJsonFile(dataPath string) ([]LinippetData, error) {
	b, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("Failed read snippet file: %w", err)
	}
	var d []LinippetData
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, err
	}
	return d, nil
}

func WriteJsonFile(dataPath string, data string) (err error) {
	d, err := ReadJsonFile(dataPath)
	if err != nil {
		// create new data when reading error
		fmt.Println(err)
		d = []LinippetData{}
	}
	file, err := os.Create(dataPath)
	if err != nil {
		return fmt.Errorf("Failed open file: %w", err)
	}
	defer func() {
		deferErr := file.Close()
		if deferErr != nil {
			err = deferErr
		}
	}()
	// Write data as JSON format file
	d = append(d, LinippetData{
		Snippet: data,
	})
	out, err := json.Marshal(&d)
	if err != nil {
		return err
	}
	if _, err := file.Write(out); err != nil {
		return err
	}
	return nil
}
