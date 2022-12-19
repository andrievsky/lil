package main

type LoopState struct {
	Continue bool
	Error    error
}

var Continue = LoopState{true, nil}
var Done = LoopState{false, nil}

type ListController struct {
	input  *Input
	view   *View
	client Client
}

func NewListController(input *Input, view *View, client Client) *ListController {
	return &ListController{input, view, client}
}

func (c *ListController) Run(path string) error {
	list, err := c.client.List(path)
	if err != nil {
		return err
	}
	c.view.List(list)

	for {
		switch c.input.PoolEvent() {
		case ResizeView:
			c.view.Resize()
			break
		case GoQuit:
			return nil
		case GoUp:
			break
		case GoDown:
			break
		case GoBack:
			break
		}
	}
}
