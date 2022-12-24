package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyItems(t *testing.T) {
	var items []Path
	model, err := NewListModel(items, 100)
	assert.Errorf(t, err, "items must contain at least one item")
	assert.Nil(t, model)
}

func TestNoVisibleItems(t *testing.T) {
	items := []Path{buildPath("abc")}
	model, err := NewListModel(items, 0)
	assert.Errorf(t, err, "maxVisibleItems must be greater than zero")
	assert.Nil(t, model)
}

func TestOneItem(t *testing.T) {
	items := []Path{buildPath("abc")}
	model, err := NewListModel(items, 1)
	assert.Nil(t, err)
	assert.Equal(t, 0, model.selectedIndex)
	assert.Equal(t, items, model.VisibleItems())
}

func TestMultipleItems(t *testing.T) {
	items := []Path{
		buildPath("a"),
		buildPath("b"),
		buildPath("c"),
	}
	model, err := NewListModel(items, 10)
	assert.Nil(t, err)
	assert.Equal(t, 0, model.selectedIndex)
	assert.Equal(t, items, model.VisibleItems())
}

func TestOnVisible(t *testing.T) {
	items := []Path{
		buildPath("a"),
		buildPath("b"),
		buildPath("c"),
	}
	model, err := NewListModel(items, 2)
	assert.Nil(t, err)
	assert.Equal(t, 0, model.selectedIndex)
	assert.Equal(t, items[:2], model.VisibleItems())
}

func TestSelectFirst(t *testing.T) {
	items := []Path{
		buildPath("a"),
		buildPath("b"),
		buildPath("c"),
	}
	model, err := NewListModel(items, 2)
	assert.Nil(t, err)
	model.Select(1)
	assert.Equal(t, items[1], model.Selected())
	assert.Equal(t, 1, model.SelectedIndex())
	assert.Equal(t, 1, model.VisibleSelectedIndex())
	assert.Equal(t, items[:2], model.VisibleItems())
}

func TestSelectLast(t *testing.T) {
	items := []Path{
		buildPath("a"),
		buildPath("b"),
		buildPath("c"),
	}
	model, err := NewListModel(items, 2)
	assert.Nil(t, err)
	model.Select(2)
	assert.Equal(t, items[2], model.Selected())
	assert.Equal(t, 2, model.SelectedIndex())
	assert.Equal(t, 1, model.VisibleSelectedIndex())
	assert.Equal(t, items[1:], model.VisibleItems())
}

func TestSelectFirstAfterLast(t *testing.T) {
	items := []Path{
		buildPath("a"),
		buildPath("b"),
		buildPath("c"),
		buildPath("d"),
	}
	model, err := NewListModel(items, 2)
	assert.Nil(t, err)
	model.Select(3)
	model.Select(1)
	assert.Equal(t, items[1], model.Selected())
	assert.Equal(t, 1, model.SelectedIndex())
	assert.Equal(t, 0, model.VisibleSelectedIndex())
	assert.Equal(t, items[1:3], model.VisibleItems())
}

func TestSelectNext(t *testing.T) {
	items := []Path{
		buildPath("a"),
		buildPath("b"),
		buildPath("c"),
	}
	model, err := NewListModel(items, 2)
	assert.Nil(t, err)
	model.SelectNext(1)
	assert.Equal(t, items[1], model.Selected())
	assert.Equal(t, 1, model.SelectedIndex())
	assert.Equal(t, 1, model.VisibleSelectedIndex())
	assert.Equal(t, items[:2], model.VisibleItems())
}

func TestSelectNextLast(t *testing.T) {
	items := []Path{
		buildPath("a"),
		buildPath("b"),
		buildPath("c"),
	}
	model, err := NewListModel(items, 2)
	assert.Nil(t, err)
	model.SelectNext(2)
	assert.Equal(t, items[2], model.Selected())
	assert.Equal(t, 2, model.SelectedIndex())
	assert.Equal(t, 1, model.VisibleSelectedIndex())
	assert.Equal(t, items[1:], model.VisibleItems())
}

func TestSelectNextAfterLast(t *testing.T) {
	items := []Path{
		buildPath("a"),
		buildPath("b"),
		buildPath("c"),
		buildPath("d"),
	}
	model, err := NewListModel(items, 2)
	assert.Nil(t, err)
	model.SelectNext(3)
	model.SelectNext(-2)
	assert.Equal(t, items[1], model.Selected())
	assert.Equal(t, 1, model.SelectedIndex())
	assert.Equal(t, 0, model.VisibleSelectedIndex())
	assert.Equal(t, items[1:3], model.VisibleItems())
}

func buildPath(path string) Path {
	return NewPath(nil, path, path, path, false)
}
