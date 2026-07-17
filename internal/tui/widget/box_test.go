package widget

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestBoxImplementsPrimitive(t *testing.T) {
	var _ Primitive = NewBox()
}

func TestBoxDrawBorderAndTitle(t *testing.T) {
	screen := newTestScreen(t)
	box := NewBox().SetBorder(true).SetTitle(" 3/10 ")
	box.SetRect(0, 0, 20, 5)
	box.Draw(screen)

	mainc, _, _, _ := screen.GetContent(0, 0)
	if mainc != tcell.RuneULCorner {
		t.Errorf("top-left = %q, want %q", mainc, tcell.RuneULCorner)
	}
	mainc, _, _, _ = screen.GetContent(19, 4)
	if mainc != tcell.RuneLRCorner {
		t.Errorf("bottom-right = %q, want %q", mainc, tcell.RuneLRCorner)
	}
	// Title is drawn on the top border starting at x+1.
	row := []rune(screenLine(screen, 0, 9))
	if got := string(row[1:8]); got != " 3/10 "+string(tcell.RuneHLine) {
		t.Errorf("title row = %q", got)
	}
}

func TestBoxInnerRect(t *testing.T) {
	box := NewBox()
	box.SetRect(2, 3, 20, 10)
	x, y, w, h := box.GetInnerRect()
	if x != 2 || y != 3 || w != 20 || h != 10 {
		t.Errorf("inner rect without border = (%d,%d,%d,%d)", x, y, w, h)
	}
	box.SetBorder(true)
	x, y, w, h = box.GetInnerRect()
	if x != 3 || y != 4 || w != 18 || h != 8 {
		t.Errorf("inner rect with border = (%d,%d,%d,%d), want (3,4,18,8)", x, y, w, h)
	}
}

func TestBoxFocus(t *testing.T) {
	box := NewBox()
	if box.HasFocus() {
		t.Error("new box should not have focus")
	}
	box.Focus()
	if !box.HasFocus() {
		t.Error("box should have focus after Focus()")
	}
	box.Blur()
	if box.HasFocus() {
		t.Error("box should not have focus after Blur()")
	}
}

func TestBoxApplyInputCapture(t *testing.T) {
	box := NewBox()
	event := tcell.NewEventKey(tcell.KeyCtrlQ, 0, tcell.ModNone)
	if got := box.ApplyInputCapture(event); got != event {
		t.Error("without capture, the event should pass through unchanged")
	}
	box.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey { return nil })
	if got := box.ApplyInputCapture(event); got != nil {
		t.Error("capture returning nil should consume the event")
	}
}
