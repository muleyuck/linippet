package linippet

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/google/uuid"
)

type Linippet struct {
	Id      string `json:"id"`
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

func AddLinippet(snippet string) error {
	linippets, err := ReadLinippets()
	if err != nil {
		// create new data when reading error
		fmt.Println(err)
		linippets = Linippets{}
	}
	linippets = append(linippets, Linippet{
		Id:      uuid.NewString(),
		Snippet: snippet,
	})
	return writeLinippets(linippets)
}

func UpdateLinippet(id, snippet string) error {
	linippets, err := ReadLinippets()
	if err != nil {
		return err
	}
	targetIndex := slices.IndexFunc(linippets, func(l Linippet) bool {
		return id == l.Id
	})
	if targetIndex == -1 {
		return fmt.Errorf("Linippet Id %s is no found", id)
	}
	linippets[targetIndex].Snippet = snippet
	return writeLinippets(linippets)
}

func RemoveLinippet(id string) error {
	linippets, err := ReadLinippets()
	if err != nil {
		return err
	}
	targetIndex := slices.IndexFunc(linippets, func(l Linippet) bool {
		return id == l.Id
	})
	if targetIndex == -1 {
		return fmt.Errorf("Linippet Id %s is no found", id)
	}
	newLinippets := slices.Delete(linippets, targetIndex, targetIndex+1)
	return writeLinippets(newLinippets)
}

func writeLinippets(linippets Linippets) (err error) {
	path, err := checkJsonPath()
	if err != nil {
		return err
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
	out, err := json.Marshal(&linippets)
	if err != nil {
		return err
	}
	if _, err := file.Write(out); err != nil {
		return err
	}
	return nil
}
