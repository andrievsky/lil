package main

import (
	"fmt"
	"math"
	"time"
)

const pageSize = 25

type Controller struct {
	input       *Input
	view        View
	client      Client
	display     *Display
	currentPath Path
}

func NewController(input *Input, view View, client Client) *Controller {
	return &Controller{input, view, client, NewDisplay(view), nil}
}

func (c *Controller) Run() error {
	if err := c.list(c.client.RootPath()); err != nil {
		return err
	}
	for {
		c.view.Render()
		switch c.input.PoolEvent() {
		case OnResize:
			c.view.RenderAll()
			break
		case GoQuit:
			return nil
		case GoUp:
			c.selectNext(-1)
			break
		case GoDown:
			c.selectNext(1)
			break
		case GoHome:
			c.selectNext(0)
			break
		case GoEnd:
			c.selectNext(math.MaxInt)
			break
		case GoPageUp:
			c.selectNext(-pageSize)
		case GoPageDown:
			c.selectNext(pageSize)
		case GoForward:
			selectedPath := c.display.list.Selected()
			if err := c.open(selectedPath); err != nil {
				return err
			}
		case GoBack:
			if err := c.back(c.currentPath); err != nil {
				return err
			}
			break
		}
	}
}

func (c *Controller) open(path Path) error {
	if path.Final() {
		err := c.content(path)
		return err
	}
	return c.list(path)
}

func (c *Controller) back(path Path) error {
	if !path.HasParent() {
		c.status("Path %s has no parent", path.GlobalPath())
		return nil
	}
	return c.open(path.Parent())
}

func (c *Controller) list(path Path) error {
	c.status("Loading...")
	c.view.Render()
	startTime := time.Now()
	list, err := c.client.List(path)
	loadingTime := time.Now().Sub(startTime)
	if err != nil {
		return err
	}
	c.currentPath = path
	c.display.list.Items(list)
	c.status("Loaded %s in %s", path.GlobalPath(), FormatEscapedTime(loadingTime))
	return nil
}

func (c *Controller) content(path Path) error {
	c.status("Loading...")
	c.view.Render()
	startTime := time.Now()
	content, err := c.client.Get(path)
	loadingTime := time.Now().Sub(startTime)
	if err != nil {
		return err
	}
	c.display.content.Set(content)
	c.status("Loaded %s in %s", path.GlobalPath(), FormatEscapedTime(loadingTime))
	return nil
}

func (c *Controller) status(msg string, a ...any) {
	c.display.status.Text(fmt.Sprintf(msg, a...))
}

func (c *Controller) selectNext(offset int) {
	c.display.list.SelectNext(offset)
	c.status("Selected %s", c.display.list.Selected().GlobalPath())
}
