package helper

import (
	"fmt"
	"regexp"
	"strings"
)

func RemoveLabelChar(text string) string {
	return regexp.MustCompile(`^>\s(.+)`).ReplaceAllString(text, "${1}")
}

func AddLabelChar(text string) string {
	return regexp.MustCompile(`^\s\s(.+)`).ReplaceAllString(text, "> ${1}")
}

func ExtractSnippetArgs(snippet string) []string {
	re := regexp.MustCompile(`\${{(\w+)}}`)
	matchArgs := re.FindAllStringSubmatch(snippet, -1)
	// without change when no args
	if len(matchArgs) <= 0 {
		return nil
	}
	linippetArgs := make([]string, 0, len(matchArgs))
	for _, matchArg := range matchArgs {
		linippetArgs = append(linippetArgs, matchArg[1])
	}
	return linippetArgs
}

func ReplaceSnippet(snippet string, index int, value string) (string, error) {
	var result string
	if index < 0 {
		return result, fmt.Errorf("Invalid index %d", index)
	}
	re := regexp.MustCompile(`(\${{\w+}})`)
	matchArgs := re.FindAllStringSubmatch(snippet, index+1)
	if len(matchArgs) <= 0 {
		return result, fmt.Errorf("Args is not found")
	}
	if len(matchArgs) <= index {
		return result, fmt.Errorf("out of args index %d", index)
	}
	return strings.Replace(snippet, matchArgs[index][1], value, 1), nil
}
