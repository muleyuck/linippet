package snippet

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	CURRENT_LABEL    = "> "
	NO_CURRENT_LABEL = "  "
)

var (
	LabelRegexp       = regexp.MustCompile(`^>\s(.+)`)
	NoLabelRegexp     = regexp.MustCompile(`^\s\s(.+)`)
	ExtractArgsRegexp = regexp.MustCompile(`\${{(\w+)}}`)
	ReplaceRegexp     = regexp.MustCompile(`(\${{\w+}})`)
)

func TrimLabel(snippet string) string {
	return LabelRegexp.ReplaceAllString(snippet, "${1}")
}

func SetNoCurrentLabel(snippet string) string {
	return NO_CURRENT_LABEL + TrimLabel(snippet)
}

func SetCurrentLabel(snippet string) string {
	return NoLabelRegexp.ReplaceAllString(snippet, CURRENT_LABEL+"${1}")
}

func AddCurrentLabel(snippet string) string {
	return CURRENT_LABEL + snippet
}

func AddNoCurrentLabel(snippet string) string {
	return NO_CURRENT_LABEL + snippet
}

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
