package list

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/style"
)

type ActionFunc func(items ...*Item) tea.Cmd

type List struct {
	Model         list.Model
	Title         string
	SelectionList bool
	ActionFunc    ActionFunc
	Hash          map[string]string
	Style         list.Styles
	id            int
	style.Frame
	*Items
}

func NewList() *List {
	m := List{
		Frame: style.DefaultFrameStyle(),
		Items: NewItems(),
	}
	m.SetAction(PrintItems)
	m.Frame.MinHeight = 10
	return &m
}

func NewListModel(w, h int, items *Items) list.Model {
	l := list.New(items.Visible(), items, w, h)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.KeyMap = ListKeyMap()
	l.Styles = ListStyles()
	return l
}
