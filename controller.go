package main

const pageSize = 10

type ListController struct {
	input    *Input
	view     *View
	client   Client
	listView *ListView
}

func NewListController(input *Input, view *View, client Client) *ListController {
	return &ListController{input, view, client, NewListView(view, 0, 0, 20, 100)}
}

func (c *ListController) Run(path string) error {
	list, err := c.client.List(path)
	if err != nil {
		return err
	}
	c.listView.Items(list)
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
		case GoBack:
			break
		}
	}
}
