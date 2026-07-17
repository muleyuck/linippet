package widget

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestModalDoneFiredByButtons(t *testing.T) {
	var gotIndex int
	var gotLabel string
	modal := NewModal().
		AddButtons([]string{"OK", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			gotIndex, gotLabel = buttonIndex, buttonLabel
		})
	modal.Focus()
	modal.HandleKey(key(tcell.KeyRight)) // OK -> Cancel
	modal.HandleKey(key(tcell.KeyEnter))
	if gotIndex != 1 || gotLabel != "Cancel" {
		t.Errorf("done = (%d, %q), want (1, Cancel)", gotIndex, gotLabel)
	}
}

func TestModalEscapeFiresDoneWithMinusOne(t *testing.T) {
	gotIndex := 99
	modal := NewModal().
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, _ string) { gotIndex = buttonIndex })
	modal.Focus()
	modal.HandleKey(key(tcell.KeyEscape))
	if gotIndex != -1 {
		t.Errorf("done index = %d, want -1", gotIndex)
	}
}

func TestModalChangedFuncReceivesInputIndex(t *testing.T) {
	var gotIndex int
	var gotValue string
	modal := NewModal().
		AddInputFields([]string{"name", "env"}, nil).
		AddButtons([]string{"OK"}).
		SetChangedFunc(func(inputIndex int, inputValue string) {
			gotIndex, gotValue = inputIndex, inputValue
		})
	modal.Focus()
	modal.HandleKey(key(tcell.KeyTab)) // -> second field
	modal.HandleKey(runeKey('x'))
	if gotIndex != 1 || gotValue != "x" {
		t.Errorf("changed = (%d, %q), want (1, x)", gotIndex, gotValue)
	}
}

func TestModalDefaultTextsAndSelectAllReplace(t *testing.T) {
	var gotValue string
	modal := NewModal().
		AddInputFields([]string{"name"}, []string{"default"}).
		AddButtons([]string{"OK"}).
		SetChangedFunc(func(_ int, inputValue string) { gotValue = inputValue })
	modal.Focus() // field gets focus -> selects all
	modal.HandleKey(runeKey('x'))
	if gotValue != "x" {
		t.Errorf("value = %q, want %q (typing replaces the default)", gotValue, "x")
	}
}

func TestModalCtrlNavigationConversions(t *testing.T) {
	modal := NewModal().
		AddInputFields([]string{"a", "b"}, nil).
		AddButtons([]string{"OK"})
	modal.Focus()
	modal.HandleKey(key(tcell.KeyCtrlN)) // -> field b (converted to Tab)
	modal.HandleKey(runeKey('z'))
	var got string
	modal.SetChangedFunc(func(_ int, v string) { got = v })
	modal.HandleKey(runeKey('!'))
	if got != "z!" {
		t.Errorf("second field text = %q, want %q", got, "z!")
	}
}

func TestModalDrawCentersItself(t *testing.T) {
	screen := newTestScreen(t) // 80x24
	modal := NewModal().
		AddButtons([]string{"OK", "Cancel"}).
		SetText("Remove?")
	modal.SetRect(0, 0, 80, 24) // the app assigns the full screen; Draw recenters
	modal.Draw(screen)

	x, y, width, height := modal.GetRect()
	// contentWidth = max(80/3, 18) = 26; outer = 30; lines = 1; form = 1; h = 7
	if width != 30 || height != 7 {
		t.Errorf("modal size = (%d,%d), want (30,7)", width, height)
	}
	if x != (80-30)/2 || y != (24-7)/2 {
		t.Errorf("modal pos = (%d,%d), want centered (25,8)", x, y)
	}
	// Border must be drawn at the modal corner.
	s, _, _ := screen.Get(x, y)
	mainc := []rune(s)[0]
	if mainc != tcell.RuneULCorner {
		t.Errorf("corner rune = %q", mainc)
	}
}
