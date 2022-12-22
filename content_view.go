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
	if c.content == nil {
		if content == nil {
			return
		}
	} else {
		if content != nil && content == c.content {
			return
		}
	}
	c.clearView()
	c.content = content
	if c.content != nil {
		c.view.DrawText(c.x, c.y, c.x+c.width-1, c.y+c.height-1, tcell.StyleDefault, c.content.Data())
	}
}

func (c *ContentView) Clear() {
	if c.content == nil {
		return
	}
	c.clearView()
}

func (c *ContentView) clearView() {
	c.view.ClearArea(c.x, c.y, c.width, c.height)
}
