package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/muleyuck/linippet/internal/fuzzy_search"
	"github.com/muleyuck/linippet/internal/linippet"
	"github.com/muleyuck/linippet/internal/snippet"
	"github.com/rivo/tview"
)

const FOCUS_LABEL = "> "

type tui struct {
	app          *tview.Application
	Result       string
	linippetArgs []string
	Submit       bool
}

type OnlyModalTui struct {
	*tui
	modal *Modal
}

func NewCreateTui() *OnlyModalTui {
	app := tview.NewApplication()
	modal := NewModal().
		AddInputFields([]string{""}, nil).
		AddTextView("").
		AddButtons([]string{"OK", "Cancel"}).
		SetText("Enter your new snippet\nYou can set argument : ${{arg_name}}")
	app.SetRoot(modal, true)
	return &OnlyModalTui{
		tui:   &tui{app: app},
		modal: modal,
	}
}

func (t *OnlyModalTui) SetAction() {
	t.modal.SetChangedFunc(func(inputIndex int, inputValue string) {
		t.Result = inputValue
		linippetArgs := snippet.ExtractSnippetArgs(inputValue)
		if len(linippetArgs) > 0 {
			text := fmt.Sprintf("This snippet have following arguments\n %v", linippetArgs)
			t.modal.textView.SetText(text)
		} else {
			t.modal.textView.SetText("")
		}
	})
	t.modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" {
			t.app.Stop()
		} else if buttonLabel == "OK" {
			t.Submit = true
			t.app.Stop()
		}
	})
	t.modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			t.app.Stop()
			return nil
		}
		return event
	})
}

func (t *OnlyModalTui) StartApp() error {
	if err := t.app.Run(); err != nil {
		t.app.Stop()
		return err
	}
	return nil
}

type listModalTui struct {
	*tui
	flex      *tview.Flex
	input     *tview.InputField
	list      *List
	linippets linippet.Linippets
	modalFunc func(string) *Modal
	SelectId  string
}

func NewRootTui() *listModalTui {
	m := newListModalTui()
	m.modalFunc = m.setRootModal
	return m
}

func NewEditTui() *listModalTui {
	m := newListModalTui()
	m.modalFunc = m.setEditModal
	return m
}

func NewRemoveTui() *listModalTui {
	m := newListModalTui()
	m.modalFunc = m.setRemoveModal
	return m
}

func newListModalTui() *listModalTui {
	app := tview.NewApplication()

	input := tview.NewInputField()
	input.SetLabel(FOCUS_LABEL)
	input.SetLabelStyle(tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault).Bold(true))
	input.SetFieldStyle(tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault))
	input.SetAcceptanceFunc(tview.InputFieldMaxLength(200)).SetFieldWidth(0)

	list := NewList()
	list.SetLabel(FOCUS_LABEL)
	list.SetBackgroundColor(tcell.ColorDefault)
	list.SetBorder(true)
	list.SetHighlightFullLine(true)
	list.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorDefault).Bold(true))
	list.SetMainTextStyle(tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault))

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(input, 1, 0, true)
	flex.AddItem(list, 0, 1, false)

	app.SetRoot(flex, true).SetFocus(input)
	return &listModalTui{
		tui:   &tui{app: app},
		flex:  flex,
		list:  list,
		input: input,
	}
}

func (t *listModalTui) SetAction() {
	t.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyCtrlN, tcell.KeyTab:
			t.offsetItem(1)
			return nil
		case tcell.KeyUp, tcell.KeyCtrlP, tcell.KeyBacktab:
			t.offsetItem(-1)
			return nil
		case tcell.KeyEnter:
			currentIndex := t.list.GetCurrentItem()
			if t.list.GetItemCount() <= currentIndex {
				t.app.Stop()
				return nil
			}
			currentText, linippetId := t.list.GetItemText(currentIndex)
			t.SelectId = linippetId
			modal := t.modalFunc(currentText)
			if modal == nil {
				t.app.Stop()
				return nil
			}
			t.input.Blur()
			t.flex.AddItem(modal, 0, 0, true)
			t.app.SetFocus(modal)
			return nil
		}
		return event
	})
	t.input.SetChangedFunc(func(text string) {
		go func() {
			t.app.QueueUpdateDraw(func() {
				t.list.Clear()
				if len(text) <= 0 {
					for _, linippet := range t.linippets {
						t.addItem(linippet.Snippet, linippet.Id, nil)
					}
				} else {
					sorted := fuzzy_search.FuzzySearch(text, t.linippets)
					for _, result := range sorted {
						t.addItem(result.Linippet.Snippet, result.Linippet.Id, result.Matches)
					}
				}
			})
		}()
	})
}

func (t *listModalTui) StartApp() error {
	if err := t.app.Run(); err != nil {
		t.app.Stop()
		return err
	}
	return nil
}

func mod(a, b int) int {
	return ((a % b) + b) % b
}

func (t *listModalTui) offsetItem(offset int) {
	currentIndex := t.list.GetCurrentItem()
	if currentIndex < 0 {
		return
	}
	itemCount := t.list.GetItemCount()
	if itemCount <= 0 {
		return
	}

	distIndex := mod(currentIndex+offset, itemCount)
	t.list.SetCurrentItem(distIndex)
}

func (t *listModalTui) addItem(mainText string, secondaryText string, matchIndices []int) {
	t.list.AddItem(mainText, secondaryText, 0, nil, matchIndices).ShowSecondaryText(false)
}

func (t *listModalTui) LazyLoadLinippet() {
	go func() {
		// time.Sleep(1 * time.Second)
		linippets, err := linippet.ReadLinippets()
		if err != nil {
			panic(err)
		}
		t.app.QueueUpdateDraw(func() {
			for _, linippet := range linippets {
				t.addItem(linippet.Snippet, linippet.Id, nil)
			}
			t.list.SetTitle(fmt.Sprintf(" %d/%d ", len(linippets), len(linippets))).SetTitleAlign(tview.AlignLeft)
		})
		t.linippets = linippets
	}()
}

func (t *listModalTui) setRootModal(currentText string) *Modal {
	linippetArgs := snippet.ExtractSnippetArgs(currentText)
	if len(linippetArgs) == 0 {
		t.Result = currentText
		return nil
	}
	t.linippetArgs = linippetArgs
	modal := NewModal().
		AddInputFields(linippetArgs, nil).
		AddButtons([]string{"OK", "Cancel"}).
		SetText(currentText)
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" {
			t.flex.RemoveItem(modal)
			t.app.SetFocus(t.input)
		} else if buttonLabel == "OK" {
			t.Result = modal.text
			t.app.Stop()
		}
	})
	modal.SetChangedFunc(func(inputIndex int, inputValue string) {
		if len(inputValue) > 0 {
			t.linippetArgs[inputIndex] = inputValue
			result, err := snippet.ReplaceSnippet(currentText, t.linippetArgs)
			if err != nil {
				return
			}
			modal.SetText(result)
		} else {
			modal.SetText(currentText)
		}
	})
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			t.flex.RemoveItem(modal)
			t.app.SetFocus(t.input)
			return nil
		}
		return event
	})

	return modal
}

func (t *listModalTui) setEditModal(currentText string) *Modal {
	modal := NewModal().
		AddInputFields([]string{""}, []string{currentText}).
		AddTextView("").
		AddButtons([]string{"OK", "Cancel"}).
		SetText("Edit snippet\nYou can set argument : ${{arg_name}}")

	t.Result = currentText
	linippetArgs := snippet.ExtractSnippetArgs(currentText)
	if len(linippetArgs) > 0 {
		text := fmt.Sprintf("This snippet have following arguments\n %v", linippetArgs)
		modal.textView.SetText(text)
	} else {
		modal.textView.SetText("")
	}
	modal.SetChangedFunc(func(inputIndex int, inputValue string) {
		t.Result = inputValue
		linippetArgs := snippet.ExtractSnippetArgs(inputValue)
		if len(linippetArgs) > 0 {
			text := fmt.Sprintf("This snippet have following arguments\n %v", linippetArgs)
			modal.textView.SetText(text)
		} else {
			modal.textView.SetText("")
		}
	})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" {
			t.flex.RemoveItem(modal)
			t.app.SetFocus(t.input)
		} else if buttonLabel == "OK" {
			t.Submit = true
			t.app.Stop()
		}
	})
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			t.flex.RemoveItem(modal)
			t.app.SetFocus(t.input)
			return nil
		}
		return event
	})

	return modal
}

func (t *listModalTui) setRemoveModal(currentText string) *Modal {
	modal := NewModal().
		AddButtons([]string{"OK", "Cancel"}).
		SetText("Remove the following snippet?\n\n" + currentText + "\n")

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" {
			t.flex.RemoveItem(modal)
			t.app.SetFocus(t.input)
		} else if buttonLabel == "OK" {
			t.Submit = true
			t.app.Stop()
		}
	})
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			t.flex.RemoveItem(modal)
			t.app.SetFocus(t.input)
			return nil
		}
		return event
	})

	return modal
}
