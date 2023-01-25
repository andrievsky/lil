package main

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

const pageSize = 25
const parallelExecutionLimit = 50

var Quit = errors.New("interrupted by user")

var EmptyPreview = NewContent(nil, "All loaded items are cached and preview automatically\n\nUse arrow keys to navigate\nEnter to open\nBackspace to go back\nQ to quit")

type Controller struct {
	input         Input
	view          View
	client        Client
	display       *Display
	currentPath   Path
	cachedList    map[Path][]Path
	cachedContent map[Path]Content
	lock          sync.RWMutex
}

func NewController(input Input, view View, client Client) *Controller {
	return &Controller{
		input,
		view,
		client,
		NewDisplay(view),
		nil,
		make(map[Path][]Path),
		make(map[Path]Content),
		sync.RWMutex{},
	}
}

func (c *Controller) Run() error {
	rootPath, err := c.client.Init()
	if err != nil {
		return err
	}
	if err = c.list(rootPath); err != nil {
		return err
	}
	for {
		c.preview(c.display.list.Selected())
		c.view.Render()
		event := c.input.PoolEvent()
		switch event {
		case OnResize:
			c.view.RenderAll()
			break
		case GoQuit:
			return Quit
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
			done, err := c.back(c.currentPath)
			if err != nil {
				return err
			}
			if done {
				return nil
			}
			break
		case OnRefresh:
			if err := c.refresh(); err != nil {
				return err
			}
			break
		default:
			{
				if event.HasKey() {
					c.display.list.SelectKey(event.Key())
				}
			}
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

func (c *Controller) back(path Path) (bool, error) {
	if !path.HasParent() {
		c.status("Path %s has no parent", path.Path())
		return true, nil
	}
	return false, c.open(path.Parent())
}

func (c *Controller) list(path Path) error {
	c.status("Loading...")
	c.view.Render()
	startTime := time.Now()
	list, err := c.loadList(path)
	if err != nil {
		return err
	}
	if err := c.preloadAllItems(list); err != nil {
		return err
	}
	loadingTime := time.Now().Sub(startTime)
	c.currentPath = path
	c.display.list.Items(list)
	c.status("Loaded %s in %s", path.Path(), FormatEscapedTime(loadingTime))
	return nil
}

func (c *Controller) content(path Path) error {
	c.status("Loading...")
	c.view.Render()
	startTime := time.Now()
	content, err := c.loadContent(path)
	if err != nil {
		return err
	}
	loadingTime := time.Now().Sub(startTime)
	c.display.content.Set(content)
	c.status("Loaded %s in %s", path.Path(), FormatEscapedTime(loadingTime))
	return nil
}

func (c *Controller) refresh() error {
	c.invalidate(c.currentPath)
	return c.list(c.currentPath)
}

func (c *Controller) invalidate(path Path) {
	list, hasList := c.cachedList[path]
	if hasList {
		for _, item := range list {
			c.invalidate(item)
		}
		delete(c.cachedList, path)
	}
	delete(c.cachedContent, path)
}

func (c *Controller) status(msg string, a ...any) {
	c.display.status.Text(fmt.Sprintf(msg, a...))
}

func (c *Controller) selectNext(offset int) {
	c.display.list.SelectNext(offset)
	selectedPath := c.display.list.Selected()
	if selectedPath != nil {
		c.status("Selected %s", c.display.list.Selected().Path())
	}
}

func (c *Controller) preview(path Path) {
	if path == nil {
		return
	}
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

func (c *Controller) loadList(path Path) ([]Path, error) {
	c.lock.RLock()
	list, cached := c.cachedList[path]
	c.lock.RUnlock()
	if !cached {
		var err error
		list, err = c.client.List(path)
		if err != nil {
			return nil, err
		}
		c.lock.Lock()
		c.cachedList[path] = list
		c.lock.Unlock()
	}
	return list, nil
}

func (c *Controller) loadContent(path Path) (Content, error) {
	c.lock.RLock()
	content, cached := c.cachedContent[path]
	c.lock.RUnlock()
	if !cached {
		var err error
		content, err = c.client.Get(path)
		if err != nil {
			return nil, err
		}
		c.lock.Lock()
		c.cachedContent[path] = content
		c.lock.Unlock()
	}
	return content, nil
}

func (c *Controller) preloadAllItems(list []Path) error {
	var tasks []func() error

	for _, path := range list {
		var p = path
		if path.Final() {
			tasks = append(tasks, func() error {
				_, err := c.loadContent(p)
				return err
			})
			continue
		}
		tasks = append(tasks, func() error {
			_, err := c.loadList(p)
			return err
		})
	}
	return ExecuteInParallel(tasks, parallelExecutionLimit, func(complete, total int) {
		c.status("Loading %d/%d", complete, total)
		c.view.Render()
	})
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
