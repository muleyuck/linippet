package widget

import "testing"

func TestVerticalLayoutDistributesRows(t *testing.T) {
	top := NewBox()
	rest := NewBox()
	layout := NewVerticalLayout().AddItem(top, 1).AddItem(rest, 0)
	layout.SetRect(0, 0, 80, 24)

	if x, y, w, h := top.GetRect(); x != 0 || y != 0 || w != 80 || h != 1 {
		t.Errorf("top rect = (%d,%d,%d,%d), want (0,0,80,1)", x, y, w, h)
	}
	if x, y, w, h := rest.GetRect(); x != 0 || y != 1 || w != 80 || h != 23 {
		t.Errorf("rest rect = (%d,%d,%d,%d), want (0,1,80,23)", x, y, w, h)
	}
}

func TestVerticalLayoutOverlayIsDrawnOnTop(t *testing.T) {
	screen := newTestScreen(t)
	base := NewBox()
	layout := NewVerticalLayout().AddItem(base, 0)
	layout.SetRect(0, 0, 80, 24)

	overlay := NewBox().SetBorder(true)
	overlay.SetRect(10, 5, 20, 5)
	layout.ShowOverlay(overlay)
	layout.Draw(screen)

	// Overlay border must be visible over the base box.
	s, _, _ := screen.Get(10, 5)
	mainc := []rune(s)[0]
	if mainc != '┌' {
		t.Errorf("overlay corner = %q, want ┌", mainc)
	}

	layout.RemoveOverlay()
	layout.Draw(screen)
	s, _, _ = screen.Get(10, 5)
	mainc = []rune(s)[0]
	if mainc != ' ' {
		t.Errorf("after RemoveOverlay, cell = %q, want space", mainc)
	}
}
