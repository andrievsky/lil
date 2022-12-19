package main

import "github.com/gdamore/tcell/v2"

type InputEvent int

const (
	GoUp InputEvent = iota
	GoDown
	GoBack
	GoQuit
	ResizeView
)

type Input struct {
	screen tcell.Screen
}

func NewInput(screen tcell.Screen) *Input {
	return &Input{screen}
}

func (i *Input) PoolEvent() InputEvent {
	for {
		event := i.screen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventResize:
			return ResizeView
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return GoQuit
			case tcell.KeyUp:
				return GoUp
			case tcell.KeyDown:
				return GoDown
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				return GoBack
			}
			switch event.Rune() {
			case 'q':
				return GoQuit
			}
		}
	}
}
