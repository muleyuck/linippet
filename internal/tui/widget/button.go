package widget

import "github.com/gdamore/tcell/v2"

// Button is a labeled button fired with Enter.
type Button struct {
	*Box
	label          string
	style          tcell.Style
	activatedStyle tcell.Style
	selected       func()
}

func NewButton(label string) *Button {
	return &Button{
		Box:            NewBox(),
		label:          label,
		style:          tcell.StyleDefault,
		activatedStyle: tcell.StyleDefault.Reverse(true),
	}
}

func (b *Button) SetSelectedFunc(handler func()) *Button {
	b.selected = handler
	return b
}

// SetStyle sets the style of the button when it is not focused.
func (b *Button) SetStyle(style tcell.Style) *Button {
	b.style = style
	return b
}

// SetActivatedStyle sets the style of the button when it is focused.
func (b *Button) SetActivatedStyle(style tcell.Style) *Button {
	b.activatedStyle = style
	return b
}

func (b *Button) GetLabel() string { return b.label }

func (b *Button) Draw(screen tcell.Screen) {
	x, y, width, _ := b.GetRect()
	style := b.style
	if b.HasFocus() {
		style = b.activatedStyle
	}
	for cx := x; cx < x+width; cx++ {
		screen.SetContent(cx, y, ' ', nil, style)
	}
	labelWidth := StringWidth(b.label)
	DrawText(screen, x+max((width-labelWidth)/2, 0), y, width, b.label, style)
}

func (b *Button) HandleKey(event *tcell.EventKey, _ func(Primitive)) {
	event = b.ApplyInputCapture(event)
	if event == nil {
		return
	}
	if event.Key() == tcell.KeyEnter {
		if b.selected != nil {
			b.selected()
		}
	}
}
