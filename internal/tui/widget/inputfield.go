package widget

import (
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// InputField is a single-line text input with emacs-style editing keys.
type InputField struct {
	*Box
	label            string
	labelStyle       tcell.Style
	labelWidth       int // minimum label column width, 0 = label's own width
	fieldStyle       tcell.Style
	text             []rune
	cursor           int // rune index of the cursor
	offset           int // first visible rune index (horizontal scroll)
	maxLength        int // maximum number of runes, 0 = unlimited
	selectAll        bool
	selectAllOnFocus bool
	changed          func(text string)
}

func NewInputField() *InputField {
	return &InputField{
		Box:        NewBox(),
		labelStyle: tcell.StyleDefault,
		fieldStyle: tcell.StyleDefault,
	}
}

func (i *InputField) SetLabel(label string) *InputField {
	i.label = label
	return i
}

func (i *InputField) SetLabelStyle(style tcell.Style) *InputField {
	i.labelStyle = style
	return i
}

// SetLabelWidth sets the minimum width of the label column, so multiple
// fields in a form can align.
func (i *InputField) SetLabelWidth(width int) *InputField {
	i.labelWidth = width
	return i
}

func (i *InputField) SetFieldStyle(style tcell.Style) *InputField {
	i.fieldStyle = style
	return i
}

func (i *InputField) SetMaxLength(maxLength int) *InputField {
	i.maxLength = maxLength
	return i
}

// SetSelectAllOnFocus makes the field select its whole text when it receives
// focus; the next typed character replaces the text.
func (i *InputField) SetSelectAllOnFocus(selectAll bool) *InputField {
	i.selectAllOnFocus = selectAll
	return i
}

func (i *InputField) SetChangedFunc(handler func(text string)) *InputField {
	i.changed = handler
	return i
}

func (i *InputField) SetText(text string) *InputField {
	i.text = []rune(text)
	if i.maxLength > 0 && len(i.text) > i.maxLength {
		i.text = i.text[:i.maxLength]
	}
	i.cursor = len(i.text)
	i.selectAll = false
	if i.changed != nil {
		i.changed(string(i.text))
	}
	return i
}

func (i *InputField) GetText() string {
	return string(i.text)
}

// Focusable implements FormItem.
func (i *InputField) Focusable() bool {
	return true
}

func (i *InputField) LabelWidth() int {
	return max(i.labelWidth, StringWidth(i.label))
}

func (i *InputField) Focus() {
	i.Box.Focus()
	if i.selectAllOnFocus && len(i.text) > 0 {
		i.selectAll = true
	}
}

func (i *InputField) Draw(screen tcell.Screen) {
	i.Box.Draw(screen)
	x, y, width, _ := i.GetInnerRect()
	if width <= 0 {
		return
	}
	labelWidth := 0
	if i.label != "" {
		labelWidth = i.LabelWidth()
		label := i.label + strings.Repeat(" ", labelWidth-StringWidth(i.label))
		DrawText(screen, x, y, min(labelWidth, width), label, i.labelStyle)
	}
	fieldX := x + labelWidth
	fieldWidth := width - labelWidth
	if fieldWidth <= 0 {
		return
	}

	style := i.fieldStyle
	if i.selectAll {
		style = style.Reverse(true)
	}

	// Fill the field background.
	for fx := fieldX; fx < fieldX+fieldWidth; fx++ {
		screen.SetContent(fx, y, ' ', nil, i.fieldStyle)
	}

	// Keep the cursor in view; one cell is reserved for the cursor itself.
	if i.cursor < i.offset {
		i.offset = i.cursor
	}
	for StringWidth(string(i.text[i.offset:i.cursor])) >= fieldWidth {
		i.offset++
	}

	DrawText(screen, fieldX, y, fieldWidth, string(i.text[i.offset:]), style)
	if i.HasFocus() {
		screen.ShowCursor(fieldX+StringWidth(string(i.text[i.offset:i.cursor])), y)
	}
}

func (i *InputField) HandleKey(event *tcell.EventKey) {
	event = i.ApplyInputCapture(event)
	if event == nil {
		return
	}

	before := string(i.text)
	switch event.Key() {
	case tcell.KeyRune:
		if i.selectAll {
			i.text = nil
			i.cursor = 0
			i.selectAll = false
		}
		if i.maxLength == 0 || len(i.text) < i.maxLength {
			i.text = slices.Insert(i.text, i.cursor, event.Rune())
			i.cursor++
		}
	case tcell.KeyLeft, tcell.KeyCtrlB:
		i.selectAll = false
		if i.cursor > 0 {
			i.cursor--
		}
	case tcell.KeyRight, tcell.KeyCtrlF:
		i.selectAll = false
		if i.cursor < len(i.text) {
			i.cursor++
		}
	case tcell.KeyHome, tcell.KeyCtrlA:
		i.selectAll = false
		i.cursor = 0
	case tcell.KeyEnd, tcell.KeyCtrlE:
		i.selectAll = false
		i.cursor = len(i.text)
	case tcell.KeyBackspace, tcell.KeyBackspace2, tcell.KeyCtrlH:
		if i.selectAll {
			i.text, i.cursor, i.selectAll = nil, 0, false
		} else if i.cursor > 0 {
			i.text = slices.Delete(i.text, i.cursor-1, i.cursor)
			i.cursor--
		}
	case tcell.KeyDelete:
		if i.selectAll {
			i.text, i.cursor, i.selectAll = nil, 0, false
		} else if i.cursor < len(i.text) {
			i.text = slices.Delete(i.text, i.cursor, i.cursor+1)
		}
	case tcell.KeyCtrlU:
		i.text, i.cursor, i.selectAll = nil, 0, false
	case tcell.KeyCtrlK:
		i.text = i.text[:i.cursor]
		i.selectAll = false
	case tcell.KeyCtrlW:
		i.selectAll = false
		start := i.cursor
		for start > 0 && i.text[start-1] == ' ' {
			start--
		}
		for start > 0 && i.text[start-1] != ' ' {
			start--
		}
		i.text = slices.Delete(i.text, start, i.cursor)
		i.cursor = start
	case tcell.KeyCtrlL:
		i.selectAll = len(i.text) > 0
	}

	if string(i.text) != before && i.changed != nil {
		i.changed(string(i.text))
	}
}
