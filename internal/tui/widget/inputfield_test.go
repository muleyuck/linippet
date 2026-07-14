package widget

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func typeKeys(input *InputField, keys ...*tcell.EventKey) {
	for _, key := range keys {
		input.HandleKey(key)
	}
}

func runeKey(r rune) *tcell.EventKey {
	return tcell.NewEventKey(tcell.KeyRune, r, tcell.ModNone)
}

func key(k tcell.Key) *tcell.EventKey {
	return tcell.NewEventKey(k, 0, tcell.ModNone)
}

func typeString(input *InputField, s string) {
	for _, r := range s {
		typeKeys(input, runeKey(r))
	}
}

func TestInputFieldTyping(t *testing.T) {
	input := NewInputField()
	typeString(input, "abc")
	if got := input.GetText(); got != "abc" {
		t.Errorf("text = %q, want %q", got, "abc")
	}
}

func TestInputFieldCursorEditing(t *testing.T) {
	input := NewInputField()
	typeString(input, "abc")
	typeKeys(input, key(tcell.KeyLeft))       // cursor between b and c
	typeKeys(input, runeKey('X'))             // abXc
	typeKeys(input, key(tcell.KeyBackspace2)) // abc
	if got := input.GetText(); got != "abc" {
		t.Errorf("text = %q, want %q", got, "abc")
	}
	typeKeys(input, key(tcell.KeyCtrlA), key(tcell.KeyDelete)) // bc
	if got := input.GetText(); got != "bc" {
		t.Errorf("text = %q, want %q", got, "bc")
	}
}

func TestInputFieldLineEditingKeys(t *testing.T) {
	input := NewInputField()
	typeString(input, "echo hello world")
	typeKeys(input, key(tcell.KeyCtrlW)) // delete word before cursor
	if got := input.GetText(); got != "echo hello " {
		t.Errorf("after Ctrl+W: text = %q, want %q", got, "echo hello ")
	}
	typeKeys(input, key(tcell.KeyCtrlU)) // clear all
	if got := input.GetText(); got != "" {
		t.Errorf("after Ctrl+U: text = %q, want empty", got)
	}
	typeString(input, "abcdef")
	typeKeys(input, key(tcell.KeyLeft), key(tcell.KeyLeft), key(tcell.KeyCtrlK))
	if got := input.GetText(); got != "abcd" {
		t.Errorf("after Ctrl+K: text = %q, want %q", got, "abcd")
	}
}

func TestInputFieldMaxLength(t *testing.T) {
	input := NewInputField().SetMaxLength(3)
	typeString(input, "abcdef")
	if got := input.GetText(); got != "abc" {
		t.Errorf("text = %q, want %q", got, "abc")
	}
}

func TestInputFieldSelectAllOnFocusReplacesTextOnType(t *testing.T) {
	input := NewInputField().SetSelectAllOnFocus(true).SetText("default")
	input.Focus()
	typeString(input, "x")
	if got := input.GetText(); got != "x" {
		t.Errorf("text = %q, want %q (typing replaces selected text)", got, "x")
	}
}

func TestInputFieldSelectAllClearedByCursorMove(t *testing.T) {
	input := NewInputField().SetSelectAllOnFocus(true).SetText("default")
	input.Focus()
	typeKeys(input, key(tcell.KeyRight))
	typeString(input, "x")
	if got := input.GetText(); got != "defaultx" {
		t.Errorf("text = %q, want %q", got, "defaultx")
	}
}

func TestInputFieldChangedFunc(t *testing.T) {
	var got []string
	input := NewInputField().SetChangedFunc(func(text string) {
		got = append(got, text)
	})
	typeString(input, "ab")
	typeKeys(input, key(tcell.KeyLeft)) // no text change: no callback
	if len(got) != 2 || got[0] != "a" || got[1] != "ab" {
		t.Errorf("changed calls = %v, want [a ab]", got)
	}
}

func TestInputFieldInputCaptureConsumesEvent(t *testing.T) {
	input := NewInputField()
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'z' {
			return nil
		}
		return event
	})
	typeString(input, "az")
	if got := input.GetText(); got != "a" {
		t.Errorf("text = %q, want %q", got, "a")
	}
}

func TestInputFieldDraw(t *testing.T) {
	screen := newTestScreen(t)
	input := NewInputField().SetLabel("> ").SetText("hello")
	input.SetRect(0, 0, 20, 1)
	input.Draw(screen)
	if got := screenLine(screen, 0, 20); got != "> hello" {
		t.Errorf("drawn = %q, want %q", got, "> hello")
	}
}

func TestInputFieldDrawScrollsToKeepCursorVisible(t *testing.T) {
	screen := newTestScreen(t)
	input := NewInputField().SetText("abcdefghij") // cursor at end
	input.SetRect(0, 0, 5, 1)
	input.Focus()
	input.Draw(screen)
	// Field width 5, one cell reserved for the cursor: "ghij" visible.
	if got := screenLine(screen, 0, 5); got != "ghij" {
		t.Errorf("drawn = %q, want %q", got, "ghij")
	}
}
