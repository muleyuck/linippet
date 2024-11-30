package tui

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/muleyuck/linippet/internal/linippet"
	"github.com/muleyuck/linippet/internal/slice"
	"github.com/muleyuck/linippet/internal/snippet"
	"github.com/rivo/tview"
)

type tui struct {
	app       *tview.Application
	flex      *tview.Flex
	input     *tview.InputField
	list      *tview.List
	linippets linippet.Linippets
	result    string
}

func NewTui() *tui {
	app := tview.NewApplication()

	input := tview.NewInputField()
	input.SetLabel(snippet.CURRENT_LABEL)
	input.SetLabelStyle(tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault).Bold(true))
	input.SetFieldStyle(tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault))
	input.SetAcceptanceFunc(tview.InputFieldMaxLength(200)).SetFieldWidth(0)

	list := tview.NewList()
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
	return &tui{
		app:   app,
		flex:  flex,
		list:  list,
		input: input,
	}
}

func (t *tui) SetAction() {
	// inputの入力イベント
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
			currentText, _ := t.list.GetItemText(currentIndex)
			linippetArgs := snippet.ExtractSnippetArgs(currentText)
			if linippetArgs == nil {
				t.result = currentText
				t.app.Stop()
				return nil
			}
			text := snippet.TrimLabel(currentText)
			argsFormModal := NewModal().
				AddInputFields(linippetArgs).
				AddButtons([]string{"OK", "Cancel"}).
				SetText(text)
			argsFormModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Cancel" {
					t.flex.RemoveItem(argsFormModal)
					t.app.SetFocus(t.input)
				} else if buttonLabel == "OK" {
					t.result = argsFormModal.text
					t.app.Stop()
				}
			})
			argsFormModal.SetChangedFunc(func(inputIndex int, inputValue string) {
				if len(inputValue) > 0 {
					result, err := snippet.ReplaceSnippet(text, inputIndex, inputValue)
					if err != nil {
						return
					}
					argsFormModal.SetText(result)
				} else {
					argsFormModal.SetText(text)
				}
			})
			argsFormModal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				switch event.Key() {
				case tcell.KeyCtrlQ:
					t.flex.RemoveItem(argsFormModal)
					t.app.SetFocus(t.input)
					return nil
				}
				return event
			})
			t.input.Blur()
			// TODO: textのLableとColorをリセット
			t.flex.AddItem(argsFormModal, 0, 0, true)
			t.app.SetFocus(argsFormModal)
			return nil
		}
		return event
	})
	// inputのChangeイベント
	t.input.SetChangedFunc(func(text string) {
		go func() {
			t.app.QueueUpdateDraw(func() {
				t.list.Clear()
				if len(text) <= 0 {
					for index, linippet := range t.linippets {
						t.addItem(index, linippet.Snippet, 0)
					}
				} else {
					// TODO: Fuzzy Search
					filtered := slice.FilterSlice(slices.Values(t.linippets), func(linippet linippet.Linippet) bool {
						return strings.Contains(linippet.Snippet, text)
					})
					sorted := slices.SortedFunc(filtered, func(a linippet.Linippet, b linippet.Linippet) int {
						return 1
					})
					currentIndex := t.list.GetCurrentItem()
					if currentIndex > len(sorted)-1 {
						currentIndex = len(sorted) - 1
					}
					for index, linippet := range sorted {
						t.addItem(index, linippet.Snippet, currentIndex)
					}
				}
			})
		}()
	})
}

func (t *tui) StartApp() error {
	if err := t.app.Run(); err != nil {
		t.app.Stop()
		return err
	}
	return nil
}

func (t *tui) GetResult() string {
	return snippet.TrimLabel(t.result)
}

func mod(a, b int) int {
	return ((a % b) + b) % b
}

func (t *tui) offsetItem(offset int) {
	currentIndex := t.list.GetCurrentItem()
	if currentIndex < 0 {
		return
	}
	itemCount := t.list.GetItemCount()
	if itemCount <= 0 {
		return
	}

	mainText, _ := t.list.GetItemText(currentIndex)
	t.list.SetItemText(currentIndex, snippet.SetNoCurretLabel(mainText), "")

	distIndex := mod(currentIndex+offset, itemCount)
	t.list.SetCurrentItem(distIndex)
	distText, _ := t.list.GetItemText(distIndex)

	t.list.SetItemText(distIndex, snippet.SetCurrentLabel(distText), "")
}

func (t *tui) addItem(nowIndex int, text string, currentIndex int) {
	var item string
	if nowIndex == currentIndex {
		item = snippet.AddCurrentLabel(text)
	} else {
		item = snippet.AddNoCurrentLabel(text)
	}
	t.list.AddItem(item, "", 0, nil).ShowSecondaryText(false)
}

func (t *tui) LazyLoadLinippet() {
	go func() {
		// time.Sleep(1 * time.Second)
		linippets, err := linippet.ReadLinippets()
		if err != nil {
			panic(err)
		}
		t.app.QueueUpdateDraw(func() {
			for index, linippet := range linippets {
				t.addItem(index, linippet.Snippet, 0)
			}
			t.list.SetTitle(fmt.Sprintf(" %d/%d ", len(linippets), len(linippets))).SetTitleAlign(tview.AlignLeft)
		})
		t.linippets = linippets
	}()
}
