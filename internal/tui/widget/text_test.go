package widget

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
)

func newTestScreen(t *testing.T) tcell.SimulationScreen {
	t.Helper()
	screen := tcell.NewSimulationScreen("")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	screen.SetSize(80, 24)
	return screen
}

// screenLine reads back a row of the simulation screen as a string
// (trailing spaces trimmed).
func screenLine(screen tcell.SimulationScreen, y, width int) string {
	var b strings.Builder
	for x := range width {
		s, _, _ := screen.Get(x, y)
		mainc := []rune(s)[0]
		b.WriteRune(mainc)
	}
	return strings.TrimRight(b.String(), " ")
}

func TestStringWidth(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"abc", 3},
		{"日本語", 6},
		{"a日b", 4},
	}
	for _, tt := range tests {
		if got := StringWidth(tt.in); got != tt.want {
			t.Errorf("StringWidth(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestWordWrap(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		width int
		want  []string
	}{
		{"fits", "hello", 10, []string{"hello"}},
		{"breaks at space", "hello wonderful world", 10, []string{"hello", "wonderful", "world"}},
		{"hard-breaks long word", "abcdefghij", 4, []string{"abcd", "efgh", "ij"}},
		{"preserves newlines", "ab\ncd", 10, []string{"ab", "cd"}},
	}
	for _, tt := range tests {
		if got := WordWrap(tt.text, tt.width); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s: WordWrap(%q, %d) = %q, want %q", tt.name, tt.text, tt.width, got, tt.want)
		}
	}
}

func TestDrawTextClipsAtMaxWidth(t *testing.T) {
	screen := newTestScreen(t)
	printed := DrawText(screen, 2, 0, 5, "hello world", tcell.StyleDefault)
	if printed != 5 {
		t.Errorf("printed width = %d, want 5", printed)
	}
	if got := screenLine(screen, 0, 80); got != "  hello" {
		t.Errorf("screen line = %q, want %q", got, "  hello")
	}
}

func TestDrawTextWideRunes(t *testing.T) {
	screen := newTestScreen(t)
	printed := DrawText(screen, 0, 0, 5, "日本語", tcell.StyleDefault)
	if printed != 4 { // "語" (width 2) does not fit into the remaining 1 cell
		t.Errorf("printed width = %d, want 4", printed)
	}
}

func TestDrawTextStyledCallsStyleAtWithByteIndex(t *testing.T) {
	screen := newTestScreen(t)
	var indices []int
	DrawTextStyled(screen, 0, 0, 80, "abc", tcell.StyleDefault,
		func(byteIndex int, base tcell.Style) tcell.Style {
			indices = append(indices, byteIndex)
			return base
		})
	if want := []int{0, 1, 2}; !reflect.DeepEqual(indices, want) {
		t.Errorf("byte indices = %v, want %v", indices, want)
	}
}
