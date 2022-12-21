package main

type ListView struct {
	view *View
	ViewPosition
	ViewSize
	labelViews    []*LabelView
	sourceItems   []Path
	selectedIndex int
}

func NewListView(view *View, x, y, width, height int) *ListView {
	return &ListView{
		view,
		ViewPosition{x, y},
		ViewSize{width, height},
		[]*LabelView{},
		[]Path{},
		0,
	}
}

func (l *ListView) Items(list []Path) {
	l.selectedIndex = 0
	l.Clear()
	views := make([]*LabelView, len(list))
	for i, item := range list {
		view := NewLabelView(l.view, item.Label(), i == l.selectedIndex, l.x, l.y+i, l.width)
		views[i] = view
	}
	l.sourceItems = list
	l.labelViews = views
}

func (l *ListView) Select(index int) {
	if index < 0 || index >= len(l.labelViews) {
		return
	}
	l.labelViews[l.selectedIndex].Select(false)
	l.selectedIndex = index
	l.labelViews[l.selectedIndex].Select(true)
}

func (l *ListView) SelectNext(step int) {
	nextIndex := l.selectedIndex + step
	if nextIndex < 0 {
		nextIndex = 0
	} else if nextIndex >= len(l.labelViews) {
		nextIndex = len(l.labelViews) - 1
	}
	l.Select(nextIndex)
}

func (l *ListView) SelectedItem() Path {
	return l.sourceItems[l.selectedIndex]
}

func (l *ListView) Clear() {
	for _, label := range l.labelViews {
		label.Clear()
	}
}
