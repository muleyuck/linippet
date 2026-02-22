package snippet

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	LabelRegexp       = regexp.MustCompile(`^>\s(.+)`)
	NoLabelRegexp     = regexp.MustCompile(`^\s\s(.+)`)
	ExtractArgsRegexp = regexp.MustCompile(`\${{(\w+)(?::([^}]*))?}}`)
	ReplaceRegexp     = regexp.MustCompile(`(\${{[^}]*}})`)
)

type Arg struct {
	Name    string
	Default string
}

func ExtractSnippetArgsWithDefaults(snippet string) []Arg {
	matchArgs := ExtractArgsRegexp.FindAllStringSubmatch(snippet, -1)
	if len(matchArgs) <= 0 {
		return nil
	}
	args := make([]Arg, 0, len(matchArgs))
	for _, matchArg := range matchArgs {
		args = append(args, Arg{Name: matchArg[1], Default: matchArg[2]})
	}
	return args
}

func ReplaceSnippet(snippet string, args []string) (string, error) {
	result := snippet
	if len(args) == 0 {
		return result, fmt.Errorf("must have args")
	}
	matchArgs := ReplaceRegexp.FindAllStringSubmatch(result, len(args))
	for index, arg := range args {
		if len(matchArgs) <= 0 {
			return result, fmt.Errorf("args is not found")
		}
		if len(matchArgs) <= index {
			return result, fmt.Errorf("out of range: index is %d", index)
		}
		result = strings.Replace(result, matchArgs[index][1], arg, 1)
	}
	return result, nil
}

// Validate snipppet is one-liner
// One-liner means that it does not contain any newline characters."
// However, line ending character strings are permitted.”
func ValidateSnippet(snippet string) error {
	if strings.ContainsAny(snippet, "\n\r") {
		return fmt.Errorf("linippet is supported only one-liner snippet")
	}
	return nil
}
