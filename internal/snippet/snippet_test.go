package snippet

import (
	"reflect"
	"testing"
)

func TestExtractSnippetArgsWithDefaults(t *testing.T) {
	tests := []struct {
		name     string
		snippet  string
		expected []Arg
	}{
		{name: "empty snippet", snippet: "", expected: nil},
		{name: "no args", snippet: "ls .", expected: nil},
		{name: "arg with default", snippet: "echo ${{greeting:hello}}", expected: []Arg{{Name: "greeting", Default: "hello"}}},
		{name: "arg without default", snippet: "echo ${{name}}", expected: []Arg{{Name: "name", Default: ""}}},
		{name: "mixed args", snippet: "echo ${{a:x}} ${{b}}", expected: []Arg{{Name: "a", Default: "x"}, {Name: "b", Default: ""}}},
		{name: "default with special chars", snippet: "curl ${{url:localhost:8080}}", expected: []Arg{{Name: "url", Default: "localhost:8080"}}},
		{name: "empty default value with colon", snippet: "echo ${{arg:}}", expected: []Arg{{Name: "arg", Default: ""}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractSnippetArgsWithDefaults(tt.snippet)
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
		{name: "out of args index", snippet: "ls ${{args}}", args: []string{}, expected: "ls ${{args}}", isOccurredError: true},
		{name: "success replace", snippet: "ls ${{args}}", args: []string{"hoge"}, expected: "ls hoge", isOccurredError: false},
		{name: "have many args", snippet: "ls ${{option}} ${{dir}}", args: []string{"hoge"}, expected: "ls hoge ${{dir}}", isOccurredError: false},
		{name: "success multiple args", snippet: "ls ${{option}} ${{dir}}", args: []string{"hoge", "fuga"}, expected: "ls hoge fuga", isOccurredError: false},
		{name: "replace arg with default syntax", snippet: "echo ${{greeting:hello}}", args: []string{"hi"}, expected: "echo hi", isOccurredError: false},
		{name: "replace using default value", snippet: "echo ${{greeting:hello}}", args: []string{"hello"}, expected: "echo hello", isOccurredError: false},
		{name: "replace mixed default and no-default", snippet: "echo ${{a:x}} ${{b}}", args: []string{"x", "y"}, expected: "echo x y", isOccurredError: false},
		{name: "replace empty placeholder", snippet: "echo ${{}}", args: []string{"hello"}, expected: "echo hello", isOccurredError: false},
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

func TestValidateSnippet(t *testing.T) {
	tests := []struct {
		name            string
		snippet         string
		isOccurredError bool
	}{
		{name: "valid one line", snippet: "echo 'hello world'", isOccurredError: false},
		{name: "valid with literal backslash-n", snippet: "echo \"Line1\\nLine2\"", isOccurredError: false},
		{name: "valid with echo -e", snippet: "echo -e \"Line1\\nLine2\"", isOccurredError: false},
		{name: "valid complex command", snippet: "ls -la | grep test && echo done", isOccurredError: false},
		{name: "invalid with actual newline", snippet: "echo \"Line1\"\necho \"Line2\"", isOccurredError: true},
		{name: "invalid with CRLF", snippet: "echo \"Line1\"\r\necho \"Line2\"", isOccurredError: true},
		{name: "invalid with CR only", snippet: "echo \"Line1\"\recho \"Line2\"", isOccurredError: true},
		{name: "invalid with multiple newlines", snippet: "echo \"Line1\"\necho \"Line2\"\necho \"Line3\"", isOccurredError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSnippet(tt.snippet)
			if (err != nil) != tt.isOccurredError {
				t.Errorf("In spite of isOccurredError = %+v, error occurred: %+v", tt.isOccurredError, err)
			}
		})
	}
}
