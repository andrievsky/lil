package main

import "time"

type ListController struct {
	list []ListItem
	view *View
}

func NewListController(list []ListItem, view *View) *ListController {
	return &ListController{list, view}
}

func (c *ListController) Run() Feedback {
	c.view.List(c.list)
	time.Sleep(5 * time.Second)
	return Feedback{
		Action: QuitAction,
	}
}
