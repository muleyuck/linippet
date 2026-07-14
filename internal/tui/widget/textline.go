package widget

import "github.com/gdamore/tcell/v2"

// TextLine is a static, non-focusable line of text, used inside a Form for
// help text.
type TextLine struct {
	*Box
	text  string
	style tcell.Style
}

func NewTextLine(text string) *TextLine {
	return &TextLine{Box: NewBox(), text: text, style: tcell.StyleDefault}
}

// Focusable implements FormItem.
func (t *TextLine) Focusable() bool {
	return false
}

func (t *TextLine) Draw(screen tcell.Screen) {
	t.Box.Draw(screen)
	x, y, width, _ := t.GetInnerRect()
	DrawText(screen, x, y, width, t.text, t.style)
}
