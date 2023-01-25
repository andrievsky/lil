package main

import "github.com/gdamore/tcell/v2"

type InputEvent int

const (
	GoUp InputEvent = iota
	GoDown
	GoHome
	GoEnd
	GoPageUp
	GoPageDown
	GoForward
	GoBack
	GoQuit
	OnResize
	OnRefresh
)

const KeyInputOffset InputEvent = 256

func (t InputEvent) HasKey() bool {
	return t >= KeyInputOffset
}

func (t InputEvent) Key() rune {
	return rune(t - KeyInputOffset)
}

func KeyInputEvent(key rune) InputEvent {
	return InputEvent(key) + KeyInputOffset
}

type Input interface {
	PoolEvent() InputEvent
}

type KeyboardInput struct {
	screen tcell.Screen
}

func NewKeyboardInput(screen tcell.Screen) Input {
	return &KeyboardInput{screen}
}

func (t *KeyboardInput) PoolEvent() InputEvent {
	for {
		event := t.screen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventResize:
			return OnResize
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return GoQuit
			case tcell.KeyUp:
				return GoUp
			case tcell.KeyDown:
				return GoDown
			case tcell.KeyHome:
				return GoHome
			case tcell.KeyEnd:
				return GoEnd
			case tcell.KeyPgUp:
				return GoPageUp
			case tcell.KeyPgDn:
				return GoPageDown
			case tcell.KeyEnter, tcell.KeyRight:
				return GoForward
			case tcell.KeyLeft, tcell.KeyBackspace, tcell.KeyBackspace2:
				return GoBack
			case tcell.KeyCtrlR:
				return OnRefresh
			}
			key := event.Rune()
			return KeyInputEvent(key)
		}
	}
}
