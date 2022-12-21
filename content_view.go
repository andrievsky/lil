package main

import "github.com/gdamore/tcell/v2"

type ContentView struct {
	view View
	ViewPosition
	ViewSize
	content Content
}

func NewContentView(view View, x, y, width, height int) *ContentView {
	return &ContentView{
		view,
		ViewPosition{x, y},
		ViewSize{width, height},
		nil,
	}
}

func (c *ContentView) Set(content Content) {
	c.content = content
	c.Clear()
	if c.content == nil {
		return
	}
	c.view.DrawText(c.x, c.y, c.x+c.width-1, c.y+c.height-1, tcell.StyleDefault, c.content.Data())
}

func (c *ContentView) Clear() {
	c.view.ClearArea(c.x, c.y, c.width, c.height)
}
