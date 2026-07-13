package widget

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func newTestForm() (*Form, *InputField, *InputField) {
	first := NewInputField().SetLabel("first")
	second := NewInputField().SetLabel("second")
	form := NewForm()
	form.AddFormItem(first)
	form.AddFormItem(NewTextLine("help")) // must be skipped in focus order
	form.AddFormItem(second)
	form.AddButton("OK", nil)
	form.AddButton("Cancel", nil)
	return form, first, second
}

func TestFormFocusStartsAtFirstItem(t *testing.T) {
	form, first, _ := newTestForm()
	form.Focus()
	if !first.HasFocus() {
		t.Error("first field should have focus")
	}
}

func TestFormTabCyclesSkippingTextLine(t *testing.T) {
	form, first, second := newTestForm()
	form.Focus()
	form.HandleKey(key(tcell.KeyTab), nil)
	if !second.HasFocus() {
		t.Fatal("Tab should move focus to the second field, skipping the TextLine")
	}
	form.HandleKey(key(tcell.KeyTab), nil) // -> OK
	form.HandleKey(key(tcell.KeyTab), nil) // -> Cancel
	form.HandleKey(key(tcell.KeyTab), nil) // wraps -> first
	if !first.HasFocus() {
		t.Error("Tab should wrap around to the first field")
	}
	form.HandleKey(key(tcell.KeyBacktab), nil) // back to Cancel
	if !form.GetButton(1).HasFocus() {
		t.Error("Backtab should wrap back to the last button")
	}
}

func TestFormEnterOnFieldMovesToNext(t *testing.T) {
	form, _, second := newTestForm()
	form.Focus()
	form.HandleKey(key(tcell.KeyEnter), nil)
	if !second.HasFocus() {
		t.Error("Enter on a field should move focus to the next target")
	}
}

func TestFormEnterOnButtonFires(t *testing.T) {
	fired := ""
	form := NewForm()
	form.AddButton("OK", func() { fired = "OK" })
	form.AddButton("Cancel", func() { fired = "Cancel" })
	form.Focus()
	form.HandleKey(key(tcell.KeyRight), nil) // OK -> Cancel
	form.HandleKey(key(tcell.KeyEnter), nil)
	if fired != "Cancel" {
		t.Errorf("fired = %q, want %q", fired, "Cancel")
	}
	form.HandleKey(key(tcell.KeyLeft), nil) // Cancel -> OK
	form.HandleKey(key(tcell.KeyEnter), nil)
	if fired != "OK" {
		t.Errorf("fired = %q, want %q", fired, "OK")
	}
}

func TestFormEscapeFiresCancel(t *testing.T) {
	canceled := false
	form, _, _ := newTestForm()
	form.SetCancelFunc(func() { canceled = true })
	form.Focus()
	form.HandleKey(key(tcell.KeyEscape), nil)
	if !canceled {
		t.Error("Escape should fire the cancel func")
	}
}

func TestFormTypingGoesToFocusedField(t *testing.T) {
	form, first, second := newTestForm()
	form.Focus()
	form.HandleKey(runeKey('a'), nil)
	form.HandleKey(key(tcell.KeyTab), nil)
	form.HandleKey(runeKey('b'), nil)
	if first.GetText() != "a" || second.GetText() != "b" {
		t.Errorf("texts = %q, %q; want %q, %q", first.GetText(), second.GetText(), "a", "b")
	}
}

func TestFormHeight(t *testing.T) {
	form, _, _ := newTestForm()
	// 3 items * 2 rows + 1 button row
	if got := form.Height(); got != 7 {
		t.Errorf("Height() = %d, want 7", got)
	}
	buttonsOnly := NewForm()
	buttonsOnly.AddButton("OK", nil)
	if got := buttonsOnly.Height(); got != 1 {
		t.Errorf("Height() = %d, want 1", got)
	}
}

func TestFormDrawCentersButtons(t *testing.T) {
	screen := newTestScreen(t)
	form := NewForm()
	form.AddButton("OK", nil)
	form.AddButton("Cancel", nil)
	form.SetRect(0, 0, 40, 1)
	form.Draw(screen)
	// Buttons row: total width = (2+4) + 2 + (6+4) = 18, centered at x=11.
	got := screenLine(screen, 0, 40)
	if got != "             OK      Cancel" {
		t.Errorf("buttons row = %q", got)
	}
}
