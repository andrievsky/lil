package main

type Display struct {
	list    *ListView
	content *ContentView
	status  *LabelView
}

func NewDisplay(view View) *Display {
	width, height := view.Size()
	borderSize := 1
	halfBodyWidth := width / 2
	bodyHeight := height - borderSize*2

	list := NewListView(view, 0, 0, halfBodyWidth-borderSize, bodyHeight)
	content := NewContentView(view, halfBodyWidth+borderSize, 0, halfBodyWidth-borderSize-1, bodyHeight)
	status := NewLabelView(view, "", false, 0, height-1, width-1)
	return &Display{
		list,
		content,
		status,
	}
}

func (l *Display) update() {

}
