package tui

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/muleyuck/linippet/internal/fuzzy_search"
	"github.com/muleyuck/linippet/internal/linippet"
	"github.com/muleyuck/linippet/internal/snippet"
	"github.com/muleyuck/linippet/internal/tui/widget"
)

const FOCUS_LABEL = "> "

type tui struct {
	app          *widget.App
	Result       string
	linippetArgs []string
	Submit       bool
}

type OnlyModalTui struct {
	*tui
	modal *widget.Modal
}

func NewCreateTui() *OnlyModalTui {
	app := widget.NewApp()
	modal := widget.NewModal().
		AddInputFields([]string{""}, nil).
		AddTextView("Syntax: ${{name}} or ${{name:default}}").
		AddButtons([]string{"OK", "Cancel"}).
		SetText("$ ")
	app.SetRoot(modal)
	return &OnlyModalTui{
		tui:   &tui{app: app},
		modal: modal,
	}
}

func (t *OnlyModalTui) SetAction() {
	t.modal.SetChangedFunc(func(inputIndex int, inputValue string) {
		t.Result = inputValue
		t.modal.SetText("$ " + snippetPreviewText(inputValue))
	})
	t.modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" || buttonIndex == -1 {
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
	t.app.SetFocus(t.modal)
	if err := t.app.Run(); err != nil {
		t.app.Stop()
		return err
	}
	return nil
}

type listModalTui struct {
	*tui
	layout       *widget.VerticalLayout
	input        *widget.InputField
	list         *widget.List
	linippets    linippet.Linippets
	modalFunc    func(string) *widget.Modal
	SelectId     string
	searchCancel context.CancelFunc
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
	app := widget.NewApp()

	input := widget.NewInputField().
		SetLabel(FOCUS_LABEL).
		SetLabelStyle(tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault).Bold(true)).
		SetFieldStyle(tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault)).
		SetMaxLength(200)

	list := widget.NewList().
		SetLabel(FOCUS_LABEL).
		SetHighlightFullLine(true).
		SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorDefault).Bold(true)).
		SetMainTextStyle(tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault))
	list.SetBackgroundColor(tcell.ColorDefault)
	list.SetBorder(true)

	layout := widget.NewVerticalLayout().
		AddItem(input, 1).
		AddItem(list, 0)

	app.SetRoot(layout)
	return &listModalTui{
		tui:    &tui{app: app},
		layout: layout,
		list:   list,
		input:  input,
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
			t.layout.ShowOverlay(modal)
			t.app.SetFocus(modal)
			return nil
		}
		return event
	})
	t.input.SetChangedFunc(func(text string) {
		if t.searchCancel != nil {
			t.searchCancel()
		}

		if len(text) <= 0 {
			t.searchCancel = nil
			t.list.Clear()
			for _, linippet := range t.linippets {
				t.addItem(linippet.Snippet, linippet.Id, nil)
			}
			t.list.SetTitle(fmt.Sprintf(" %d/%d ", len(t.linippets), len(t.linippets)))
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		t.searchCancel = cancel
		go func() {
			sorted := fuzzy_search.FuzzySearch(ctx, text, t.linippets)
			if sorted == nil {
				return
			}
			t.app.QueueUpdateDraw(func() {
				if t.input.GetText() != text {
					return
				}
				t.list.Clear()
				for _, result := range sorted {
					t.addItem(result.Linippet.Snippet, result.Linippet.Id, result.Matches)
				}
				t.list.SetTitle(fmt.Sprintf(" %d/%d ", len(sorted), len(t.linippets)))
			})
		}()
	})
}

func (t *listModalTui) StartApp() error {
	t.app.SetFocus(t.input)
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
	t.list.AddItem(mainText, secondaryText, matchIndices)
}

func (t *listModalTui) LazyLoadLinippet() {
	go func() {
		linippets, err := linippet.ReadLinippets()
		if err != nil {
			panic(err)
		}
		t.app.QueueUpdateDraw(func() {
			for _, linippet := range linippets {
				t.addItem(linippet.Snippet, linippet.Id, nil)
			}
			t.list.SetTitle(fmt.Sprintf(" %d/%d ", len(linippets), len(linippets)))
		})
		t.linippets = linippets
	}()
}

func (t *listModalTui) closeModal() {
	t.layout.RemoveOverlay()
	t.app.SetFocus(t.input)
}

func (t *listModalTui) setRootModal(currentText string) *widget.Modal {
	args := snippet.ExtractSnippetArgsWithDefaults(currentText)
	if len(args) == 0 {
		t.Result = currentText
		return nil
	}
	argNames := make([]string, len(args))
	t.linippetArgs = make([]string, len(args))
	for i, arg := range args {
		argNames[i] = arg.Name
		t.linippetArgs[i] = arg.Default
	}
	modal := widget.NewModal().
		AddInputFields(argNames, t.linippetArgs).
		AddButtons([]string{"OK", "Cancel"}).
		SetText("$ " + snippetPreviewText(currentText))
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" || buttonIndex == -1 {
			t.closeModal()
		} else if buttonLabel == "OK" {
			result, err := snippet.ReplaceSnippet(currentText, t.linippetArgs)
			if err != nil {
				result = currentText
			}
			t.Result = result
			t.app.Stop()
		}
	})
	modal.SetChangedFunc(func(inputIndex int, inputValue string) {
		t.linippetArgs[inputIndex] = inputValue
		result, err := snippet.ReplaceSnippet(currentText, t.linippetArgs)
		if err != nil {
			modal.SetText("$ " + currentText)
			return
		}
		modal.SetText("$ " + result)
	})
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			t.closeModal()
			return nil
		}
		return event
	})

	return modal
}

func snippetPreviewText(snippetText string) string {
	args := snippet.ExtractSnippetArgsWithDefaults(snippetText)
	if len(args) == 0 {
		return snippetText
	}
	previewArgs := make([]string, len(args))
	for i, arg := range args {
		if arg.Default != "" {
			previewArgs[i] = arg.Default
		} else {
			previewArgs[i] = "<" + arg.Name + ">"
		}
	}
	result, err := snippet.ReplaceSnippet(snippetText, previewArgs)
	if err != nil {
		return snippetText
	}
	return result
}

func (t *listModalTui) setEditModal(currentText string) *widget.Modal {
	modal := widget.NewModal().
		AddInputFields([]string{""}, []string{currentText}).
		AddTextView("Syntax: ${{name}} or ${{name:default}}").
		AddButtons([]string{"OK", "Cancel"}).
		SetText("$ " + snippetPreviewText(currentText))

	t.Result = currentText
	modal.SetChangedFunc(func(inputIndex int, inputValue string) {
		t.Result = inputValue
		modal.SetText("$ " + snippetPreviewText(inputValue))
	})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" || buttonIndex == -1 {
			t.closeModal()
		} else if buttonLabel == "OK" {
			t.Submit = true
			t.app.Stop()
		}
	})
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			t.closeModal()
			return nil
		}
		return event
	})

	return modal
}

func (t *listModalTui) setRemoveModal(currentText string) *widget.Modal {
	modal := widget.NewModal().
		AddButtons([]string{"OK", "Cancel"}).
		SetText("Remove the following snippet?\n\n" + currentText + "\n")

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" || buttonIndex == -1 {
			t.closeModal()
		} else if buttonLabel == "OK" {
			t.Submit = true
			t.app.Stop()
		}
	})
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			t.closeModal()
			return nil
		}
		return event
	})

	return modal
}
