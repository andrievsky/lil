package main

import (
	"errors"
	"fmt"
)

const pageSize = 10

type ListController struct {
	input       *Input
	view        *View
	client      Client
	listView    *ListView
	currentPath Path
}

func NewListController(input *Input, view *View, client Client) *ListController {
	return &ListController{input, view, client, NewListView(view, 0, 0, 40, 100), nil}
}

func (c *ListController) Run() error {
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
			c.listView.SelectNext(-1)
			break
		case GoDown:
			c.listView.SelectNext(1)
			break
		case GoHome:
			c.listView.Select(0)
			break
		case GoEnd:
			c.listView.Select(len(c.listView.labelViews) - 1)
			break
		case GoPageUp:
			c.listView.SelectNext(-pageSize)
		case GoPageDown:
			c.listView.SelectNext(pageSize)
		case GoForward:
			selectedItem := c.listView.SelectedItem()
			if err := c.list(selectedItem); err != nil {
				return err
			}
		case GoBack:
			if !c.currentPath.HasParent() {
				return errors.New(fmt.Sprintf("path %s has no parent", c.currentPath))
			}
			if err := c.list(c.currentPath.Parent()); err != nil {
				return err
			}
			break
		}
	}
}

func (c *ListController) list(path Path) error {
	list, err := c.client.List(path)
	if err != nil {
		return err
	}
	c.currentPath = path
	c.listView.Items(list)
	return nil
}
