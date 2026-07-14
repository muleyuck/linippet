package widget

import "github.com/gdamore/tcell/v2"

// FormItem is an element that can be laid out in a Form.
type FormItem interface {
	Primitive
	// Focusable reports whether the item participates in focus cycling.
	Focusable() bool
}

// Form lays out form items vertically with a centered button row at the
// bottom, and cycles keyboard focus between them.
type Form struct {
	*Box
	items                []FormItem
	buttons              []*Button
	focusedIndex         int // index into focusTargets()
	cancel               func()
	buttonStyle          tcell.Style
	buttonActivatedStyle tcell.Style
}

func NewForm() *Form {
	return &Form{
		Box:                  NewBox(),
		buttonStyle:          tcell.StyleDefault,
		buttonActivatedStyle: tcell.StyleDefault.Reverse(true),
	}
}

func (f *Form) AddFormItem(item FormItem) *Form {
	f.items = append(f.items, item)
	return f
}

func (f *Form) AddButton(label string, selected func()) *Form {
	button := NewButton(label).
		SetSelectedFunc(selected).
		SetStyle(f.buttonStyle).
		SetActivatedStyle(f.buttonActivatedStyle)
	f.buttons = append(f.buttons, button)
	return f
}

func (f *Form) GetButton(index int) *Button { return f.buttons[index] }

func (f *Form) GetButtonCount() int { return len(f.buttons) }

// SetCancelFunc sets the handler fired when the user presses Escape.
func (f *Form) SetCancelFunc(cancel func()) *Form {
	f.cancel = cancel
	return f
}

// SetButtonStyle sets the style applied to buttons added afterwards.
func (f *Form) SetButtonStyle(style tcell.Style) *Form {
	f.buttonStyle = style
	return f
}

// SetButtonActivatedStyle sets the focused style applied to buttons added
// afterwards.
func (f *Form) SetButtonActivatedStyle(style tcell.Style) *Form {
	f.buttonActivatedStyle = style
	return f
}

// Height returns the number of rows the form occupies: one row per item plus
// a blank row after each, and one row for the buttons.
func (f *Form) Height() int {
	height := len(f.items) * 2
	if len(f.buttons) > 0 {
		height++
	}
	return height
}

// focusTargets returns the focusable elements in cycling order: focusable
// items first, then buttons.
func (f *Form) focusTargets() []Primitive {
	targets := make([]Primitive, 0, len(f.items)+len(f.buttons))
	for _, item := range f.items {
		if item.Focusable() {
			targets = append(targets, item)
		}
	}
	for _, button := range f.buttons {
		targets = append(targets, button)
	}
	return targets
}

func (f *Form) Focus() {
	f.Box.Focus()
	targets := f.focusTargets()
	if len(targets) == 0 {
		return
	}
	if f.focusedIndex >= len(targets) {
		f.focusedIndex = 0
	}
	for index, target := range targets {
		if index == f.focusedIndex {
			target.Focus()
		} else {
			target.Blur()
		}
	}
}

func (f *Form) Blur() {
	f.Box.Blur()
	for _, target := range f.focusTargets() {
		target.Blur()
	}
}

func (f *Form) shiftFocus(offset int) {
	targets := f.focusTargets()
	count := len(targets)
	if count == 0 {
		return
	}
	targets[f.focusedIndex].Blur()
	f.focusedIndex = ((f.focusedIndex+offset)%count + count) % count
	targets[f.focusedIndex].Focus()
}

func (f *Form) HandleKey(event *tcell.EventKey) {
	event = f.ApplyInputCapture(event)
	if event == nil {
		return
	}
	targets := f.focusTargets()
	if len(targets) == 0 {
		return
	}
	current := targets[f.focusedIndex]
	_, onButton := current.(*Button)

	switch event.Key() {
	case tcell.KeyTab:
		f.shiftFocus(1)
		return
	case tcell.KeyBacktab:
		f.shiftFocus(-1)
		return
	case tcell.KeyEscape:
		if f.cancel != nil {
			f.cancel()
		}
		return
	case tcell.KeyEnter:
		if !onButton {
			f.shiftFocus(1)
			return
		}
	case tcell.KeyRight:
		if onButton {
			f.shiftFocus(1)
			return
		}
	case tcell.KeyLeft:
		if onButton {
			f.shiftFocus(-1)
			return
		}
	}
	current.HandleKey(event)
}

func (f *Form) Draw(screen tcell.Screen) {
	f.Box.Draw(screen)
	x, y, width, _ := f.GetInnerRect()

	// Align field labels.
	maxLabelWidth := 0
	for _, item := range f.items {
		if input, ok := item.(*InputField); ok {
			maxLabelWidth = max(maxLabelWidth, input.LabelWidth())
		}
	}

	row := y
	for _, item := range f.items {
		if input, ok := item.(*InputField); ok {
			input.SetLabelWidth(maxLabelWidth)
		}
		item.SetRect(x, row, width, 1)
		item.Draw(screen)
		row += 2
	}

	if len(f.buttons) == 0 {
		return
	}
	buttonsWidth := 0
	for _, button := range f.buttons {
		buttonsWidth += StringWidth(button.GetLabel()) + 4 + 2
	}
	buttonsWidth -= 2
	buttonX := x + max((width-buttonsWidth)/2, 0)
	for _, button := range f.buttons {
		buttonWidth := StringWidth(button.GetLabel()) + 4
		button.SetRect(buttonX, row, buttonWidth, 1)
		button.Draw(screen)
		buttonX += buttonWidth + 2
	}
}
