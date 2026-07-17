package tui

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/muleyuck/linippet/internal/linippet"
)

func newTestScreen(t *testing.T) tcell.SimulationScreen {
	t.Helper()
	screen := tcell.NewSimulationScreen("")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	screen.SetSize(80, 24)
	return screen
}

func setTestLinippets(target *listModalTui, linippets linippet.Linippets) {
	target.linippets = linippets
	for _, item := range linippets {
		target.addItem(item.Snippet, item.Id, nil)
	}
	target.list.SetTitle(" test ")
}

// waitFor polls a condition on the event-loop goroutine via QueueUpdateDraw.
func waitFor(t *testing.T, target *listModalTui, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		result := make(chan bool, 1)
		target.app.QueueUpdateDraw(func() { result <- condition() })
		select {
		case ok := <-result:
			if ok {
				return
			}
		case <-time.After(100 * time.Millisecond):
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("condition not met before deadline")
}

func typeText(screen tcell.SimulationScreen, text string) {
	for _, r := range text {
		screen.InjectKey(tcell.KeyRune, r, tcell.ModNone)
	}
}

func TestRootTuiSnippetWithoutArgsReturnsImmediately(t *testing.T) {
	target := NewRootTui()
	screen := newTestScreen(t)
	target.app.SetScreen(screen)
	target.SetAction()
	setTestLinippets(target, linippet.Linippets{
		{Id: "id-1", Snippet: "ls -la"},
	})
	done := make(chan error, 1)
	go func() { done <- target.StartApp() }()

	screen.InjectKey(tcell.KeyEnter, 0, tcell.ModNone)

	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if target.Result != "ls -la" {
		t.Errorf("Result = %q, want %q", target.Result, "ls -la")
	}
	if target.SelectId != "id-1" {
		t.Errorf("SelectId = %q, want %q", target.SelectId, "id-1")
	}
}

func TestRootTuiSnippetWithArgsOpensModalAndReplaces(t *testing.T) {
	target := NewRootTui()
	screen := newTestScreen(t)
	target.app.SetScreen(screen)
	target.SetAction()
	setTestLinippets(target, linippet.Linippets{
		{Id: "id-1", Snippet: "echo ${{name}}"},
	})
	done := make(chan error, 1)
	go func() { done <- target.StartApp() }()

	screen.InjectKey(tcell.KeyEnter, 0, tcell.ModNone) // open the arg modal
	typeText(screen, "world")                          // fill the arg
	screen.InjectKey(tcell.KeyEnter, 0, tcell.ModNone) // field -> OK button
	screen.InjectKey(tcell.KeyEnter, 0, tcell.ModNone) // press OK

	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if target.Result != "echo world" {
		t.Errorf("Result = %q, want %q", target.Result, "echo world")
	}
}

func TestRootTuiCtrlQClosesModalAndReturnsToInput(t *testing.T) {
	target := NewRootTui()
	screen := newTestScreen(t)
	target.app.SetScreen(screen)
	target.SetAction()
	setTestLinippets(target, linippet.Linippets{
		{Id: "id-1", Snippet: "echo ${{name}}"},
	})
	done := make(chan error, 1)
	go func() { done <- target.StartApp() }()

	screen.InjectKey(tcell.KeyEnter, 0, tcell.ModNone) // open modal
	screen.InjectKey(tcell.KeyCtrlQ, 0, tcell.ModNone) // close modal
	// Typing must reach the root input again -> fuzzy filter runs.
	typeText(screen, "zzz")
	waitFor(t, target, func() bool { return target.list.GetItemCount() == 0 })

	target.app.QueueUpdateDraw(func() { target.app.Stop() })
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestRootTuiFuzzyFilterNarrowsList(t *testing.T) {
	target := NewRootTui()
	screen := newTestScreen(t)
	target.app.SetScreen(screen)
	target.SetAction()
	setTestLinippets(target, linippet.Linippets{
		{Id: "id-1", Snippet: "echo hello"},
		{Id: "id-2", Snippet: "ls -la"},
	})
	done := make(chan error, 1)
	go func() { done <- target.StartApp() }()

	typeText(screen, "ls")
	waitFor(t, target, func() bool {
		if target.list.GetItemCount() != 1 {
			return false
		}
		main, _ := target.list.GetItemText(0)
		return main == "ls -la"
	})

	screen.InjectKey(tcell.KeyEnter, 0, tcell.ModNone)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if target.Result != "ls -la" {
		t.Errorf("Result = %q, want %q", target.Result, "ls -la")
	}
}

func TestRootTuiArrowKeysMoveSelection(t *testing.T) {
	target := NewRootTui()
	screen := newTestScreen(t)
	target.app.SetScreen(screen)
	target.SetAction()
	setTestLinippets(target, linippet.Linippets{
		{Id: "id-1", Snippet: "first"},
		{Id: "id-2", Snippet: "second"},
	})
	done := make(chan error, 1)
	go func() { done <- target.StartApp() }()

	screen.InjectKey(tcell.KeyDown, 0, tcell.ModNone)
	waitFor(t, target, func() bool { return target.list.GetCurrentItem() == 1 })
	screen.InjectKey(tcell.KeyDown, 0, tcell.ModNone) // wraps to 0
	waitFor(t, target, func() bool { return target.list.GetCurrentItem() == 0 })

	target.app.QueueUpdateDraw(func() { target.app.Stop() })
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestCreateTuiSubmit(t *testing.T) {
	target := NewCreateTui()
	screen := newTestScreen(t)
	target.app.SetScreen(screen)
	target.SetAction()
	done := make(chan error, 1)
	go func() { done <- target.StartApp() }()

	typeText(screen, "echo hi")
	screen.InjectKey(tcell.KeyEnter, 0, tcell.ModNone) // field -> OK
	screen.InjectKey(tcell.KeyEnter, 0, tcell.ModNone) // press OK

	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if !target.Submit {
		t.Error("Submit should be true after OK")
	}
	if target.Result != "echo hi" {
		t.Errorf("Result = %q, want %q", target.Result, "echo hi")
	}
}

func TestCreateTuiCtrlQQuitsWithoutSubmit(t *testing.T) {
	target := NewCreateTui()
	screen := newTestScreen(t)
	target.app.SetScreen(screen)
	target.SetAction()
	done := make(chan error, 1)
	go func() { done <- target.StartApp() }()

	screen.InjectKey(tcell.KeyCtrlQ, 0, tcell.ModNone)

	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if target.Submit {
		t.Error("Submit should be false after Ctrl+Q")
	}
}

func TestRemoveTuiOkSubmits(t *testing.T) {
	target := NewRemoveTui()
	screen := newTestScreen(t)
	target.app.SetScreen(screen)
	target.SetAction()
	setTestLinippets(target, linippet.Linippets{
		{Id: "id-1", Snippet: "dangerous command"},
	})
	done := make(chan error, 1)
	go func() { done <- target.StartApp() }()

	screen.InjectKey(tcell.KeyEnter, 0, tcell.ModNone) // open confirm modal
	screen.InjectKey(tcell.KeyEnter, 0, tcell.ModNone) // press OK (first button)

	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if !target.Submit || target.SelectId != "id-1" {
		t.Errorf("Submit = %v, SelectId = %q; want true, id-1", target.Submit, target.SelectId)
	}
}
