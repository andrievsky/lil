package main

type ListView struct {
	view View
	ViewPosition
	ViewSize
	views []*LabelView
	model *ListModel
}

func NewListView(view View, x, y, width, height int) *ListView {
	return &ListView{
		view,
		ViewPosition{x, y},
		ViewSize{width, height},
		buildViews(view, x, y, width, height),
		nil,
	}
}

func (l *ListView) Items(list []Path) error {
	l.Clear()
	model, err := NewListModel(list, l.height)
	if err != nil {
		return err
	}
	l.model = model
	l.sync()
	return nil
}

func (l *ListView) Select(index int) {
	l.model.Select(index)
	l.sync()
}

func (l *ListView) SelectNext(step int) {
	l.model.SelectNext(step)
	l.sync()
}

func (l *ListView) Selected() Path {
	return l.model.Selected()
}

func (l *ListView) Clear() {
	for _, view := range l.views {
		view.Clear()
	}
}

func buildViews(view View, x, y, width, height int) []*LabelView {
	views := make([]*LabelView, height)
	for i := 0; i < height; i++ {
		views[i] = NewLabelView(view, "", false, x, y+i, width)
	}
	return views
}

func (l *ListView) sync() {
	if l.model == nil {
		l.Clear()
		return
	}
	visibleItems := l.model.VisibleItems()
	selectedIndex := l.model.VisibleSelectedIndex()
	max := Min(len(l.views), len(visibleItems))
	for i := 0; i < max; i++ {
		l.views[i].TextAndSelect(visibleItems[i].Label(), i == selectedIndex)
	}
}
