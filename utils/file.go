package utils

import (
	"bufio"
	"errors"
	"os"
	"path"
	"path/filepath"
)

const LINIPPET_DIR = ".linippet"
const SNIPPET_FILE_NAME = "snippets.txt"

func GetSnippetFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("can not access directory: " + homeDir)
	}
	return filepath.Join(homeDir, LINIPPET_DIR, SNIPPET_FILE_NAME), nil
}

func checkDirectory(p string) error {
	baseDir := path.Dir(p)
	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}
	// make dir when target dir is not found
	return os.MkdirAll(baseDir, 0755)
}

func ReadFileLines(p string) ([]string, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, errors.New("Failed open snippet file: " + err.Error())
	}
	defer f.Close()
	var res []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func WriteFile(p string, data string) (err error) {
	err = checkDirectory(p)
	if err != nil {
		return errors.New(path.Dir(p) + " is not found: " + err.Error())
	}
	file, err := os.OpenFile(p, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return errors.New("Failed open file: " + err.Error())
	}
	defer func() {
		deferErr := file.Close()
		if deferErr != nil {
			err = deferErr
		}
	}()
	_, err = file.WriteString(data)
	if err != nil {
		return errors.New("Failed write command line: " + err.Error())
	}
	return nil
}
