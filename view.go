package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type View struct {
	screen tcell.Screen
}

func NewView() (*View, error) {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("can't create presenter screen: %w", err)
	}
	if err = screen.Init(); err != nil {
		return nil, fmt.Errorf("can't create presenter screen: %w", err)
	}
	screen.Clear()
	return &View{screen}, nil
}

func (v *View) List(list []ListItem) {
	v.screen.Clear()
	style := tcell.StyleDefault
	for i, item := range list {
		emitStr(v.screen, 0, i, style, item.Label)
	}
	v.screen.Show()
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func (v *View) Close() {
	v.screen.Fini()
}
