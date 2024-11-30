package helper

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
		name           string
		snippet        string
		index          int
		value          string
		expected       string
		isOccuredError bool
	}{
		{name: "empty snippet", snippet: "", index: 0, value: "", expected: "", isOccuredError: true},
		{name: "not args", snippet: "ls .", index: 0, value: "hoge", expected: "", isOccuredError: true},
		{name: "invalid arg character", snippet: "ls ${args}", index: 0, value: "hoge", expected: "", isOccuredError: true},
		{name: "out of args index", snippet: "ls ${{args}}", index: 1, value: "hoge", expected: "", isOccuredError: true},
		{name: "success replace", snippet: "ls ${{args}}", index: 0, value: "hoge", expected: "ls hoge", isOccuredError: false},
		{name: "have many args", snippet: "ls ${{option}} ${{dir}}", index: 1, value: "hoge", expected: "ls ${{option}} hoge", isOccuredError: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ReplaceSnippet(tt.snippet, tt.index, tt.value)
			if err != nil != tt.isOccuredError {
				t.Errorf("unexpected error: %+v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("result is %+v, but expected is %+v", result, tt.expected)
			}
		})
	}
}
