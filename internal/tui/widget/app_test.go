package widget

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
)

// keyRecorder records the keys it receives.
type keyRecorder struct {
	*Box
	keys []tcell.Key
}

func (r *keyRecorder) HandleKey(event *tcell.EventKey, _ func(Primitive)) {
	r.keys = append(r.keys, event.Key())
}

func startApp(t *testing.T, root Primitive) (*App, tcell.SimulationScreen, chan error) {
	t.Helper()
	screen := newTestScreen(t)
	app := NewApp()
	app.SetScreen(screen)
	app.SetRoot(root)
	done := make(chan error, 1)
	go func() { done <- app.Run() }()
	return app, screen, done
}

func TestAppDispatchesKeysToFocusedPrimitive(t *testing.T) {
	recorder := &keyRecorder{Box: NewBox()}
	screen := newTestScreen(t)
	app := NewApp()
	app.SetScreen(screen)
	app.SetRoot(recorder)
	// SetFocus must happen before Run starts (or via QueueUpdateDraw).
	app.SetFocus(recorder)
	done := make(chan error, 1)
	go func() { done <- app.Run() }()

	screen.InjectKey(tcell.KeyRune, 'a', tcell.ModNone)

	// Wait until the key has been recorded before stopping. A Stop queued
	// right after InjectKey would race the events and updates channels:
	// select does not order two simultaneously-ready cases, so Stop could
	// win before the key is dispatched. Polling via QueueUpdateDraw reads
	// recorder.keys on the event-loop goroutine, so it stays race-free.
	deadline := time.Now().Add(time.Second)
	for {
		recorded := make(chan bool, 1)
		app.QueueUpdateDraw(func() { recorded <- len(recorder.keys) > 0 })
		if <-recorded {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("key was never recorded")
		}
	}
	app.QueueUpdateDraw(func() { app.Stop() })
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if len(recorder.keys) != 1 || recorder.keys[0] != tcell.KeyRune {
		t.Errorf("recorded keys = %v, want [KeyRune]", recorder.keys)
	}
}

func TestAppQueueUpdateDrawBeforeRun(t *testing.T) {
	screen := newTestScreen(t)
	app := NewApp()
	app.SetScreen(screen)
	app.SetRoot(NewBox())

	ran := make(chan struct{})
	// Queued before Run starts: must not be dropped.
	app.QueueUpdateDraw(func() { close(ran) })

	done := make(chan error, 1)
	go func() { done <- app.Run() }()

	select {
	case <-ran:
	case <-time.After(time.Second):
		t.Fatal("queued update was not delivered")
	}
	app.QueueUpdateDraw(func() { app.Stop() })
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestAppStopIsIdempotent(t *testing.T) {
	app, _, done := startApp(t, NewBox())
	app.QueueUpdateDraw(func() {
		app.Stop()
		app.Stop()
	})
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	// QueueUpdateDraw after Stop must not block.
	finished := make(chan struct{})
	go func() {
		app.QueueUpdateDraw(func() {})
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("QueueUpdateDraw blocked after Stop")
	}
}
