package main

import "github.com/gdamore/tcell/v2"

const truncatedLabelSuffix = "..."

type LabelView struct {
	view        *View
	displayText []rune
	text        string
	selected    bool
	ViewPosition
	ViewSize
}

func NewLabelView(view *View, text string, selected bool, x, y, maxWidth int) *LabelView {
	label := &LabelView{
		view,
		formatText(text, maxWidth),
		text,
		selected,
		ViewPosition{x, y},
		ViewSize{maxWidth, 1},
	}
	label.update()
	return label

}

func (l *LabelView) Text(text string) {
	l.text = text
	l.displayText = formatText(text, l.width)
	l.update()
}

func (l *LabelView) Select(selected bool) {
	if selected == l.selected {
		return
	}
	l.selected = selected
	l.update()
}

func (l *LabelView) update() {
	l.view.DrawLabel(l.x, l.y, l.ViewSize.width, l.style(), l.displayText)
}

func (l *LabelView) style() tcell.Style {
	if l.selected {
		return tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	}
	return tcell.StyleDefault
}

func formatText(text string, maxWidth int) []rune {
	if len(text) <= maxWidth {
		return []rune(text)
	}
	return append([]rune(text)[:maxWidth-len(truncatedLabelSuffix)], []rune(truncatedLabelSuffix)...)
}
