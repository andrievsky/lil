package main

import (
	"github.com/gdamore/tcell/v2"
)

type View interface {
	RenderAll()
	Render()
	Size() (int, int)
	Clear()
	DrawLabel(x, y, maxWidth int, style tcell.Style, text []rune)
	ClearArea(x, y, width, height int)
	DrawText(x1, y1, x2, y2 int, style tcell.Style, text string)
}

type ViewImpl struct {
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

func NewView(screen tcell.Screen) View {
	return &ViewImpl{screen: screen}
}

func (v *ViewImpl) RenderAll() {
	v.screen.Sync()
}

func (v *ViewImpl) Render() {
	v.screen.Show()
}

func (v *ViewImpl) Size() (int, int) {
	return v.screen.Size()
}

func (v *ViewImpl) Clear() {
	v.screen.Clear()
}

func (v *ViewImpl) DrawLabel(x, y, maxWidth int, style tcell.Style, text []rune) {
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

func (v *ViewImpl) ClearArea(x, y, width, height int) {
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			v.screen.SetContent(x+i, y+j, ' ', nil, tcell.StyleDefault)
		}
	}
}

func (v *ViewImpl) DrawText(x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		if r == '\n' {
			row++
			col = x1
			continue
		}
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
