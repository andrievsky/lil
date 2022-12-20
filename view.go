package main

import (
	"github.com/gdamore/tcell/v2"
)

type View struct {
	screen tcell.Screen
}

type ViewPosition struct {
	x int
	y int
}

type ViewSize struct {
	width  int
	height int
}

func NewView(screen tcell.Screen) *View {
	return &View{screen: screen}
}

func (v *View) RenderAll() {
	v.screen.Sync()
}

func (v *View) Render() {
	v.screen.Show()
}

func (v *View) Clear() {
	v.screen.Clear()
}

func (v *View) DrawLabel(x, y, maxWidth int, style tcell.Style, text []rune) {
	if len(text) > maxWidth {
		text = text[:maxWidth]
	}
	for i, r := range text {
		v.screen.SetContent(x+i, y, r, nil, style)
	}
	for i := len(text); i < maxWidth; i++ {
		v.screen.SetContent(x+i, y, ' ', nil, style)
	}
}

func (v *View) DrawText(x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		v.screen.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}
