package linippet

import (
	"encoding/json"
	"fmt"
	"os"
)

type Linippet struct {
	Snippet string `json:"snippet"`
}
type Linippets []Linippet

func ReadLinippets() (Linippets, error) {
	path, err := checkJsonPath()
	if err != nil {
		return nil, err
	}
	return readJson(path)
}

func readJson(path string) (Linippets, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed read snippet file: %w", err)
	}
	var linippets Linippets
	if err := json.Unmarshal(b, &linippets); err != nil {
		return nil, err
	}
	return linippets, nil
}

func WriteLinippets(snippet string) (err error) {
	path, err := checkJsonPath()
	if err != nil {
		return err
	}
	linippets, err := readJson(path)
	if err != nil {
		// create new data when reading error
		fmt.Println(err)
		linippets = Linippets{}
	}
	file, err := os.Create(path)
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
	linippets = append(linippets, Linippet{
		Snippet: snippet,
	})
	out, err := json.Marshal(&linippets)
	if err != nil {
		return err
	}
	if _, err := file.Write(out); err != nil {
		return err
	}
	return nil
}
