package snippet

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	LabelRegexp       = regexp.MustCompile(`^>\s(.+)`)
	NoLabelRegexp     = regexp.MustCompile(`^\s\s(.+)`)
	ExtractArgsRegexp = regexp.MustCompile(`\${{(\w+)}}`)
	ReplaceRegexp     = regexp.MustCompile(`(\${{\w+}})`)
)

func ExtractSnippetArgs(snippet string) []string {
	matchArgs := ExtractArgsRegexp.FindAllStringSubmatch(snippet, -1)
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
	matchArgs := ReplaceRegexp.FindAllStringSubmatch(snippet, index+1)
	if len(matchArgs) <= 0 {
		return result, fmt.Errorf("Args is not found")
	}
	if len(matchArgs) <= index {
		return result, fmt.Errorf("Out of range: index is %d", index)
	}
	return strings.Replace(snippet, matchArgs[index][1], value, 1), nil
}
