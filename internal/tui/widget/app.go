package widget

import (
	"sync"

	"github.com/gdamore/tcell/v2"
)

// App owns the tcell screen, runs the event loop, and manages focus.
type App struct {
	mu       sync.Mutex
	screen   tcell.Screen
	root     Primitive
	focus    Primitive
	updates  chan func()
	done     chan struct{}
	stopOnce sync.Once
}

func NewApp() *App {
	return &App{
		updates: make(chan func(), 64),
		done:    make(chan struct{}),
	}
}

// SetScreen injects an already-initialized screen. Used by tests with
// tcell.SimulationScreen; when set, Run does not create or Init a screen.
func (a *App) SetScreen(screen tcell.Screen) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.screen = screen
}

func (a *App) SetRoot(p Primitive) *App {
	a.root = p
	return a
}

// SetFocus moves keyboard focus to p. Must be called from the event-loop
// goroutine (i.e. from key handlers or queued updates) once Run has started.
func (a *App) SetFocus(p Primitive) {
	if a.focus != nil {
		a.focus.Blur()
	}
	a.focus = p
	if p != nil {
		p.Focus()
	}
}

// QueueUpdateDraw schedules update to run on the event-loop goroutine,
// followed by a redraw. Safe to call from any goroutine. Updates queued
// after Stop are dropped.
func (a *App) QueueUpdateDraw(update func()) {
	select {
	case a.updates <- update:
	case <-a.done:
	}
}

// Stop ends the event loop. It is idempotent. Call it from key handlers or
// queued updates.
func (a *App) Stop() {
	a.stopOnce.Do(func() {
		close(a.done)
		a.mu.Lock()
		screen := a.screen
		a.mu.Unlock()
		if screen != nil {
			screen.Fini()
		}
	})
}

// Run creates the screen if none was injected and runs the event loop until
// Stop is called.
func (a *App) Run() error {
	a.mu.Lock()
	screen := a.screen
	a.mu.Unlock()
	if screen == nil {
		var err error
		screen, err = tcell.NewScreen()
		if err != nil {
			return err
		}
		if err = screen.Init(); err != nil {
			return err
		}
		a.mu.Lock()
		a.screen = screen
		a.mu.Unlock()
	}

	// Restore the terminal on panics in event handlers.
	defer func() {
		if r := recover(); r != nil {
			a.Stop()
			panic(r)
		}
	}()

	events := make(chan tcell.Event, 16)
	go func() {
		for {
			event := screen.PollEvent()
			if event == nil {
				close(events)
				return
			}
			events <- event
		}
	}()

	for {
		a.draw(screen)
		select {
		case event, ok := <-events:
			if !ok {
				return nil
			}
			switch event := event.(type) {
			case *tcell.EventKey:
				if event.Key() == tcell.KeyCtrlC {
					a.Stop()
					break
				}
				if a.focus != nil {
					a.focus.HandleKey(event)
				}
			case *tcell.EventResize:
				screen.Sync()
			}
		case update := <-a.updates:
			update()
		}
		select {
		case <-a.done:
			return nil
		default:
		}
	}
}

func (a *App) draw(screen tcell.Screen) {
	if a.root == nil {
		return
	}
	width, height := screen.Size()
	screen.HideCursor()
	a.root.SetRect(0, 0, width, height)
	a.root.Draw(screen)
	screen.Show()
}
