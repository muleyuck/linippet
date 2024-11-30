package tui

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/muleyuck/linippet/internal/file"
	"github.com/muleyuck/linippet/internal/helper"
	"github.com/rivo/tview"
)

type tui struct {
	app       *tview.Application
	flex      *tview.Flex
	input     *tview.InputField
	list      *tview.List
	linippets []file.LinippetData
	result    string
}

func NewTui() *tui {
	app := tview.NewApplication()

	input := tview.NewInputField()
	input.SetLabel("> ")
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
			currentText, _ := t.list.GetItemText(currentIndex)
			linippetArgs := helper.ExtractSnippetArgs(currentText)
			if linippetArgs == nil {
				t.result = currentText
				t.app.Stop()
			}
			argsFormModal := NewModal().
				AddInputFields(linippetArgs).
				AddButtons([]string{"OK", "Cancel"}).
				SetText(helper.RemoveLabelChar(currentText))
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
					result, err := helper.ReplaceSnippet(currentText, inputIndex, inputValue)
					if err != nil {
						return
					}
					argsFormModal.SetText(result)
				} else {
					argsFormModal.SetText(currentText)
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
		}
		return event
	})
	// inputのChangeイベント
	t.input.SetChangedFunc(func(text string) {
		go func() {
			t.app.QueueUpdateDraw(func() {
				t.list.Clear()
				if len(text) <= 0 {
					for _, linippet := range t.linippets {
						t.addItem(linippet.Snippet)
					}
				} else {
					filtered := helper.FilterSlice(slices.Values(t.linippets), func(linippet file.LinippetData) bool {
						return strings.Contains(linippet.Snippet, text)
					})
					sorted := slices.SortedFunc(filtered, func(a file.LinippetData, b file.LinippetData) int {
						return 1
					})
					currentIndex := t.list.GetCurrentItem()
					if currentIndex > len(sorted)-1 {
						currentIndex = len(sorted) - 1
					}
					for index, linippet := range sorted {
						if currentIndex == index {
							t.addItem("> " + linippet.Snippet)
						} else {
							t.addItem("  " + linippet.Snippet)
						}
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
	return helper.RemoveLabelChar(t.result)
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
	// 移動前のItemのLableを削除
	mainText, _ := t.list.GetItemText(currentIndex)
	removed := fmt.Sprintf("  %s", helper.RemoveLabelChar(mainText))
	t.list.SetItemText(currentIndex, removed, "")
	// 移動
	distIndex := mod(currentIndex+offset, itemCount)
	t.list.SetCurrentItem(distIndex)
	// 移動先のLabelをSet
	distText, _ := t.list.GetItemText(distIndex)

	t.list.SetItemText(distIndex, helper.AddLabelChar(distText), "")
}

func (t *tui) addItem(item string) {
	t.list.AddItem(item, "", 0, nil).ShowSecondaryText(false)
}

func (t *tui) LazyLoadLinippet() {
	go func() {
		time.Sleep(1 * time.Second)
		dataPath, err := file.CheckDataPath()
		if err != nil {
			panic(err)
		}
		linippets, err := file.ReadJsonFile(dataPath)
		if err != nil {
			panic(err)
		}
		t.app.QueueUpdateDraw(func() {
			for index, linippet := range linippets {
				if index == 0 {
					t.addItem("> " + linippet.Snippet)
				} else {
					t.addItem("  " + linippet.Snippet)
				}
			}
			t.list.SetTitle(fmt.Sprintf(" %d/%d ", len(linippets), len(linippets))).SetTitleAlign(tview.AlignLeft)
		})
		t.linippets = linippets
	}()
}
