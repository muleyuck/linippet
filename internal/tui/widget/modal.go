package widget

import "github.com/gdamore/tcell/v2"

// Modal is a centered window with a message text, optional input fields, and
// a button row. It positions and sizes itself on Draw.
type Modal struct {
	*Box
	form      *Form
	text      string
	textStyle tcell.Style
	changed   func(inputIndex int, inputValue string)
	done      func(buttonIndex int, buttonLabel string)
}

func NewModal() *Modal {
	m := &Modal{
		Box:       NewBox(),
		textStyle: tcell.StyleDefault,
	}
	m.SetBorder(true)
	m.form = NewForm().
		SetButtonStyle(tcell.StyleDefault).
		SetButtonActivatedStyle(tcell.StyleDefault.Background(tcell.ColorGray).Bold(true))
	m.form.SetCancelFunc(func() {
		if m.done != nil {
			m.done(-1, "")
		}
	})
	return m
}

// AddInputFields adds one input field per label. texts, when non-nil,
// provides initial values. Fields select their whole text on focus.
func (m *Modal) AddInputFields(labels []string, texts []string) *Modal {
	for index, label := range labels {
		text := ""
		if index < len(texts) {
			text = texts[index]
		}
		input := NewInputField().
			SetLabel(label).
			SetLabelStyle(tcell.StyleDefault.Foreground(tcell.ColorYellow)).
			SetFieldStyle(tcell.StyleDefault.Background(tcell.ColorGray)).
			SetText(text).
			SetSelectAllOnFocus(true).
			SetChangedFunc(func(value string) {
				if m.changed != nil {
					m.changed(index, value)
				}
			})
		m.form.AddFormItem(input)
	}
	m.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyCtrlN:
			return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
		case tcell.KeyUp, tcell.KeyCtrlP:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
		case tcell.KeyCtrlF:
			return tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone)
		case tcell.KeyCtrlB:
			return tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone)
		}
		return event
	})
	return m
}

// AddTextView adds a static, non-focusable help text row.
func (m *Modal) AddTextView(text string) *Modal {
	m.form.AddFormItem(NewTextLine(text))
	return m
}

// AddButtons adds one button per label. Pressing a button fires the done
// handler with the button's index and label.
func (m *Modal) AddButtons(labels []string) *Modal {
	for index, label := range labels {
		m.form.AddButton(label, func() {
			if m.done != nil {
				m.done(index, label)
			}
		})
	}
	return m
}

// SetText sets the message text shown above the form. It may contain line
// breaks and is word-wrapped to the modal width.
func (m *Modal) SetText(text string) *Modal {
	m.text = text
	return m
}

func (m *Modal) SetChangedFunc(handler func(inputIndex int, inputValue string)) *Modal {
	m.changed = handler
	return m
}

// SetDoneFunc sets the handler fired when a button is pressed. On Escape it
// is fired with index -1 and an empty label.
func (m *Modal) SetDoneFunc(handler func(buttonIndex int, buttonLabel string)) *Modal {
	m.done = handler
	return m
}

func (m *Modal) Focus() { m.form.Focus() }

func (m *Modal) Blur() { m.form.Blur() }

func (m *Modal) HasFocus() bool { return m.form.HasFocus() }

func (m *Modal) HandleKey(event *tcell.EventKey, setFocus func(Primitive)) {
	event = m.ApplyInputCapture(event)
	if event == nil {
		return
	}
	m.form.HandleKey(event, setFocus)
}

func (m *Modal) Draw(screen tcell.Screen) {
	screenWidth, screenHeight := screen.Size()

	// Width: at least a third of the screen, wide enough for the buttons.
	buttonsWidth := 0
	for i := range m.form.GetButtonCount() {
		buttonsWidth += StringWidth(m.form.GetButton(i).GetLabel()) + 4 + 2
	}
	buttonsWidth -= 2
	contentWidth := max(screenWidth/3, buttonsWidth)

	lines := WordWrap(m.text, contentWidth-2)
	// 2 border rows + 1 top padding + text + 1 blank + form + 1 bottom padding.
	height := min(len(lines)+m.form.Height()+5, screenHeight)
	width := contentWidth + 4
	m.SetRect((screenWidth-width)/2, (screenHeight-height)/2, width, height)

	m.Box.Draw(screen)
	x, y, innerWidth, _ := m.GetInnerRect()
	textY := y + 1
	for i, line := range lines {
		DrawText(screen, x+2, textY+i, innerWidth-4, line, m.textStyle)
	}
	m.form.SetRect(x+2, textY+len(lines)+1, innerWidth-4, m.form.Height())
	m.form.Draw(screen)
}
