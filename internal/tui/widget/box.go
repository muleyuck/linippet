package widget

import "github.com/gdamore/tcell/v2"

// Box is the base building block: a rectangle with an optional border and
// title. Widgets embed Box to inherit geometry, focus state, input capture,
// and border drawing.
type Box struct {
	x, y, width, height int
	border              bool
	title               string
	backgroundColor     tcell.Color
	focused             bool
	inputCapture        func(event *tcell.EventKey) *tcell.EventKey
}

func NewBox() *Box {
	return &Box{backgroundColor: tcell.ColorDefault}
}

func (b *Box) SetRect(x, y, width, height int) {
	b.x, b.y, b.width, b.height = x, y, width, height
}

func (b *Box) GetRect() (int, int, int, int) {
	return b.x, b.y, b.width, b.height
}

// GetInnerRect returns the drawable area inside the border.
func (b *Box) GetInnerRect() (int, int, int, int) {
	if !b.border {
		return b.x, b.y, b.width, b.height
	}
	return b.x + 1, b.y + 1, b.width - 2, b.height - 2
}

func (b *Box) SetBorder(border bool) *Box {
	b.border = border
	return b
}

// SetTitle sets the text drawn left-aligned on the top border.
func (b *Box) SetTitle(title string) *Box {
	b.title = title
	return b
}

func (b *Box) SetBackgroundColor(color tcell.Color) *Box {
	b.backgroundColor = color
	return b
}

// SetInputCapture installs a function that intercepts key events before the
// widget handles them. Returning nil consumes the event; returning a
// different event replaces it.
func (b *Box) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *Box {
	b.inputCapture = capture
	return b
}

// ApplyInputCapture runs the installed input capture, if any.
func (b *Box) ApplyInputCapture(event *tcell.EventKey) *tcell.EventKey {
	if b.inputCapture != nil {
		return b.inputCapture(event)
	}
	return event
}

func (b *Box) Draw(screen tcell.Screen) {
	if b.width <= 0 || b.height <= 0 {
		return
	}
	style := tcell.StyleDefault.Background(b.backgroundColor)
	for y := b.y; y < b.y+b.height; y++ {
		for x := b.x; x < b.x+b.width; x++ {
			screen.SetContent(x, y, ' ', nil, style)
		}
	}
	if b.border && b.width >= 2 && b.height >= 2 {
		for x := b.x + 1; x < b.x+b.width-1; x++ {
			screen.SetContent(x, b.y, tcell.RuneHLine, nil, style)
			screen.SetContent(x, b.y+b.height-1, tcell.RuneHLine, nil, style)
		}
		for y := b.y + 1; y < b.y+b.height-1; y++ {
			screen.SetContent(b.x, y, tcell.RuneVLine, nil, style)
			screen.SetContent(b.x+b.width-1, y, tcell.RuneVLine, nil, style)
		}
		screen.SetContent(b.x, b.y, tcell.RuneULCorner, nil, style)
		screen.SetContent(b.x+b.width-1, b.y, tcell.RuneURCorner, nil, style)
		screen.SetContent(b.x, b.y+b.height-1, tcell.RuneLLCorner, nil, style)
		screen.SetContent(b.x+b.width-1, b.y+b.height-1, tcell.RuneLRCorner, nil, style)
		if b.title != "" {
			DrawText(screen, b.x+1, b.y, b.width-2, b.title, style)
		}
	}
}

func (b *Box) HandleKey(_ *tcell.EventKey) {}

func (b *Box) Focus() {
	b.focused = true
}

func (b *Box) Blur() {
	b.focused = false
}

func (b *Box) HasFocus() bool {
	return b.focused
}
