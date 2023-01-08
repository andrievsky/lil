package main

import (
	"fmt"
	"math"
	"os"
)

type ClientController struct {
	input   Input
	view    View
	display *Display
	path    string
}

func NewClientController(input Input, view View, path string) *ClientController {
	return &ClientController{
		input,
		view,
		NewDisplay(view),
		path,
	}
}

func (c *ClientController) Pick() (Client, error) {
	clientList, err := listClients(c.path)
	if err != nil {
		return nil, err
	}
	c.display.list.Items(clientList)
	if err != nil {
		return nil, err
	}
	for {
		c.view.Render()
		event := c.input.PoolEvent()
		switch event {
		case OnResize:
			c.view.RenderAll()
			break
		case GoQuit:
			return nil, nil
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
			return NewCmdClient(selectedPath.Label(), selectedPath.Path())
		default:
			{
				if event.HasKey() {
					c.display.list.SelectKey(event.Key())
				}
			}
		}
	}
}

func (c *ClientController) selectNext(offset int) {
	c.display.list.SelectNext(offset)
	selectedPath := c.display.list.Selected()
	if selectedPath != nil {
		c.display.status.Text(fmt.Sprintf("Selected %s", c.display.list.Selected().Path()))
	}
}

func listClients(path string) ([]Path, error) {
	var clientList []Path
	fileInfo, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range fileInfo {
		if !file.IsDir() {
			continue
		}
		clientList = append(clientList, NewPath(nil, path+"/"+file.Name()+"/", file.Name(), false))
	}

	return clientList, nil
}
