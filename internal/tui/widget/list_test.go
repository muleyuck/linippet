package widget

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestListAddClearAndItemAccess(t *testing.T) {
	list := NewList()
	list.AddItem("echo hello", "id-1", nil)
	list.AddItem("ls -la", "id-2", nil)
	if got := list.GetItemCount(); got != 2 {
		t.Fatalf("count = %d, want 2", got)
	}
	main, secondary := list.GetItemText(1)
	if main != "ls -la" || secondary != "id-2" {
		t.Errorf("item 1 = (%q, %q)", main, secondary)
	}
	list.Clear()
	if got := list.GetItemCount(); got != 0 {
		t.Errorf("count after Clear = %d, want 0", got)
	}
	if got := list.GetCurrentItem(); got != 0 {
		t.Errorf("current after Clear = %d, want 0", got)
	}
}

func TestListSetCurrentItemClamps(t *testing.T) {
	list := NewList()
	list.AddItem("a", "", nil)
	list.AddItem("b", "", nil)
	list.SetCurrentItem(5)
	if got := list.GetCurrentItem(); got != 1 {
		t.Errorf("current = %d, want 1 (clamped)", got)
	}
	list.SetCurrentItem(-1)
	if got := list.GetCurrentItem(); got != 0 {
		t.Errorf("current = %d, want 0 (clamped)", got)
	}
}

func TestListDrawSelectedLabelAndItems(t *testing.T) {
	screen := newTestScreen(t)
	list := NewList().SetLabel("> ")
	list.AddItem("first", "", nil)
	list.AddItem("second", "", nil)
	list.SetRect(0, 0, 40, 10)
	list.SetCurrentItem(1)
	list.Draw(screen)

	if got := screenLine(screen, 0, 40); got != "  first" {
		t.Errorf("row 0 = %q, want %q", got, "  first")
	}
	if got := screenLine(screen, 1, 40); got != "> second" {
		t.Errorf("row 1 = %q, want %q", got, "> second")
	}
}

func TestListDrawHighlightsMatchedBytes(t *testing.T) {
	screen := newTestScreen(t)
	list := NewList()
	list.AddItem("abc", "", []int{0, 2}) // highlight a and c
	list.SetRect(0, 0, 40, 10)
	list.Draw(screen)

	_, _, styleA, _ := screen.GetContent(0, 0)
	fgA, _, _ := styleA.Decompose()
	if fgA != tcell.ColorGreen {
		t.Errorf("matched cell fg = %v, want green", fgA)
	}
	_, _, styleB, _ := screen.GetContent(1, 0)
	fgB, _, _ := styleB.Decompose()
	if fgB == tcell.ColorGreen {
		t.Error("unmatched cell must not be green")
	}
}

func TestListDrawScrollsToKeepCurrentInView(t *testing.T) {
	screen := newTestScreen(t)
	list := NewList()
	for _, s := range []string{"a", "b", "c", "d", "e"} {
		list.AddItem(s, "", nil)
	}
	list.SetRect(0, 0, 40, 3) // only 3 rows visible
	list.SetCurrentItem(4)
	list.Draw(screen)
	// Items c, d, e visible; e on the last row.
	if got := screenLine(screen, 2, 40); got != "e" {
		t.Errorf("last row = %q, want %q", got, "e")
	}
	if got := screenLine(screen, 0, 40); got != "c" {
		t.Errorf("first row = %q, want %q", got, "c")
	}
}

func TestListHighlightFullLine(t *testing.T) {
	screen := newTestScreen(t)
	selectedStyle := tcell.StyleDefault.Background(tcell.ColorGray)
	list := NewList().SetHighlightFullLine(true).SetSelectedStyle(selectedStyle)
	list.AddItem("x", "", nil)
	list.SetRect(0, 0, 10, 3)
	list.Draw(screen)
	_, _, style, _ := screen.GetContent(9, 0) // rightmost cell of the row
	_, bg, _ := style.Decompose()
	if bg != tcell.ColorGray {
		t.Errorf("full-line highlight bg = %v, want gray", bg)
	}
}
