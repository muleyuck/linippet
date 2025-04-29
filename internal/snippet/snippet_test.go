package snippet

import (
	"reflect"
	"testing"
)

func TestGetSnippetArgs(t *testing.T) {
	tests := []struct {
		name     string
		snippet  string
		expected []string
	}{
		{name: "empty snippet", snippet: "", expected: nil},
		{name: "not args", snippet: "ls .", expected: nil},
		{name: "invalid arg character", snippet: "ls ${args}", expected: nil},
		{name: "too many arg character", snippet: "ls ${{{args}}}", expected: nil},
		{name: "has one arg", snippet: "ls ${{option}}", expected: []string{"option"}},
		{name: "invalid second args", snippet: "ls ${{option}} ${{{dir}}} ", expected: []string{"option"}},
		{name: "have many args", snippet: "ls ${{option}} ${{dir}} ", expected: []string{"option", "dir"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractSnippetArgs(tt.snippet)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("result is %+v, but expected is %+v", result, tt.expected)
			}
		})
	}
}

func TestReplaceSnippet(t *testing.T) {
	tests := []struct {
		name            string
		snippet         string
		args            []string
		expected        string
		isOccurredError bool
	}{
		{name: "empty snippet", snippet: "", args: []string{}, expected: "", isOccurredError: true},
		{name: "not args", snippet: "ls .", args: []string{"hoge"}, expected: "ls .", isOccurredError: true},
		{name: "invalid arg character", snippet: "ls ${args}", args: []string{"hoge"}, expected: "ls ${args}", isOccurredError: true},
		{name: "out of args index", snippet: "ls ${{args}}", args: []string{}, expected: "ls ${args}}", isOccurredError: true},
		{name: "success replace", snippet: "ls ${{args}}", args: []string{"hoge"}, expected: "ls hoge", isOccurredError: false},
		{name: "have many args", snippet: "ls ${{option}} ${{dir}}", args: []string{"hoge"}, expected: "ls hoge ${{dir}}", isOccurredError: false},
		{name: "success multiple args", snippet: "ls ${{option}} ${{dir}}", args: []string{"hoge", "fuga"}, expected: "ls hoge fuga", isOccurredError: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ReplaceSnippet(tt.snippet, tt.args)
			if err != nil != tt.isOccurredError {
				t.Errorf("unexpected error: %+v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("result is %+v, but expected is %+v", result, tt.expected)
			}
		})
	}
}
