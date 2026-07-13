// Package widget is a minimal keyboard-only widget toolkit built directly on
// tcell. It intentionally supports no mouse handling and no style-tag markup.
package widget

import "github.com/gdamore/tcell/v2"

// Primitive is a UI element that can draw itself and handle key events.
type Primitive interface {
	// Draw renders the primitive onto the screen.
	Draw(screen tcell.Screen)
	// SetRect sets the position and size of the primitive.
	SetRect(x, y, width, height int)
	// GetRect returns the position and size of the primitive.
	GetRect() (x, y, width, height int)
	// HandleKey processes a key event. Implementations may call setFocus to
	// move the application focus to another primitive.
	HandleKey(event *tcell.EventKey, setFocus func(Primitive))
	// Focus marks the primitive as having keyboard focus.
	Focus()
	// Blur removes keyboard focus from the primitive.
	Blur()
	// HasFocus reports whether the primitive has keyboard focus.
	HasFocus() bool
}
