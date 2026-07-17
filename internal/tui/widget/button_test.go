package widget

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestButtonEnterFiresSelected(t *testing.T) {
	fired := false
	button := NewButton("OK").SetSelectedFunc(func() { fired = true })
	button.HandleKey(key(tcell.KeyEnter))
	if !fired {
		t.Error("Enter should fire the selected func")
	}
}

func TestButtonDrawCentersLabel(t *testing.T) {
	screen := newTestScreen(t)
	button := NewButton("OK")
	button.SetRect(0, 0, 6, 1) // label width 2 + padding 4
	button.Draw(screen)
	if got := screenLine(screen, 0, 6); got != "  OK" {
		t.Errorf("drawn = %q, want %q (2 leading cells)", got, "  OK")
	}
}

func TestTextLineIsNotFocusable(t *testing.T) {
	line := NewTextLine("help text")
	if line.Focusable() {
		t.Error("TextLine must not be focusable")
	}
}

func TestTextLineDraw(t *testing.T) {
	screen := newTestScreen(t)
	line := NewTextLine("Syntax: ${{name}}")
	line.SetRect(0, 0, 40, 1)
	line.Draw(screen)
	if got := screenLine(screen, 0, 40); got != "Syntax: ${{name}}" {
		t.Errorf("drawn = %q", got)
	}
}
