package main

import (
	"fmt"
	"math"
	"time"
)

const pageSize = 25

var EmptyPreview = NewContent(nil, "All loaded items are cached and preview automatically\n\nUse arrow keys to navigate\nEnter to open\nBackspace to go back\nQ to quit")

type Controller struct {
	input         *Input
	view          View
	client        Client
	display       *Display
	currentPath   Path
	cachedList    map[Path][]Path
	cachedContent map[Path]Content
}

func NewController(input *Input, view View, client Client) *Controller {
	return &Controller{
		input,
		view,
		client,
		NewDisplay(view),
		nil,
		make(map[Path][]Path),
		make(map[Path]Content)}
}

func (c *Controller) Run() error {
	if err := c.list(c.client.RootPath()); err != nil {
		return err
	}
	for {
		c.preview(c.display.list.Selected())
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
			c.selectNext(0)
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
	list, cached := c.cachedList[path]
	if !cached {
		var err error
		list, err = c.client.List(path)
		if err != nil {
			return err
		}
		c.cachedList[path] = list
	}
	loadingTime := time.Now().Sub(startTime)
	c.currentPath = path
	c.display.list.Items(list)
	c.status("Loaded %s %s in %s", path.GlobalPath(), optionalText(cached, "from cache", ""), FormatEscapedTime(loadingTime))
	return nil
}

func (c *Controller) content(path Path) error {
	c.status("Loading...")
	c.view.Render()
	startTime := time.Now()
	content, cached := c.cachedContent[path]
	if !cached {
		var err error
		content, err = c.client.Get(path)
		if err != nil {
			return err
		}
		c.cachedContent[path] = content
	}
	loadingTime := time.Now().Sub(startTime)
	c.display.content.Set(content)
	c.status("Loaded %s %s in %s", path.GlobalPath(), optionalText(cached, "from cache", ""), FormatEscapedTime(loadingTime))
	return nil
}

func (c *Controller) status(msg string, a ...any) {
	c.display.status.Text(fmt.Sprintf(msg, a...))
}

func (c *Controller) selectNext(offset int) {
	c.display.list.SelectNext(offset)
	selectedPath := c.display.list.Selected()
	if selectedPath != nil {
		c.status("Selected %s", c.display.list.Selected().GlobalPath())
	}
}

func (c *Controller) preview(path Path) {
	if content := c.cachedContent[path]; content != nil {
		c.display.content.Set(content)
		return
	}
	if list := c.cachedList[path]; list != nil {
		c.display.content.Set(listContent(path, list))
		return
	}
	c.display.content.Set(EmptyPreview)
	return
}

func listContent(path Path, list []Path) Content {
	return NewContent(path, fmt.Sprintf("List of %v items", len(list)))
}

func optionalText(defined bool, some, none string) string {
	if defined {
		return some
	}
	return none
}
