package main

import "errors"

type ListModel struct {
	Items             []Path
	maxVisibleItems   int
	firstVisibleIndex int
	selectedIndex     int
}

func NewListModel(items []Path, maxVisibleItems int) (*ListModel, error) {
	if len(items) < 1 {
		return nil, errors.New("items must contain at least one item")
	}
	if maxVisibleItems < 1 {
		return nil, errors.New("maxVisibleItems must be greater than zero")
	}
	return &ListModel{
		items,
		maxVisibleItems,
		0,
		0,
	}, nil
}

func (l *ListModel) Select(index int) {
	if index < 0 {
		index = 0
	} else if index >= len(l.Items) {
		index = len(l.Items) - 1
	}
	if index == l.selectedIndex {
		return
	}
	l.selectedIndex = index

	if index < l.firstVisibleIndex {
		l.firstVisibleIndex = index
		return
	}
	if index > l.lastVisibleIndex() {
		l.firstVisibleIndex = index - l.maxVisibleItems + 1
	}
}

func (l *ListModel) SelectNext(offset int) {
	l.Select(l.selectedIndex + offset)
}

func (l *ListModel) Selected() Path {
	return l.Items[l.selectedIndex]
}

func (l *ListModel) SelectedIndex() int {
	return l.selectedIndex
}

func (l *ListModel) VisibleSelectedIndex() int {
	return l.selectedIndex - l.firstVisibleIndex
}

func (l *ListModel) VisibleItems() []Path {
	last := l.lastVisibleIndex()
	return l.Items[l.firstVisibleIndex : last+1]
}

func (l *ListModel) lastVisibleIndex() int {
	index := l.firstVisibleIndex + l.maxVisibleItems - 1
	if index >= len(l.Items) {
		index = len(l.Items) - 1
	}
	return index
}
