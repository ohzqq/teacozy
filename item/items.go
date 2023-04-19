package item

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
)

type List struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	items []Item
}

type Props struct {
	Choices    Choices
	Selectable bool
	Selected   map[int]struct{}
	Cursor     int
	Start      int
	End        int
	Search     string
}

func NewList() *List {
	return &List{}
}

func (c *List) Init(props Props) tea.Cmd {
	c.UpdateProps(props)

	items := props.Choices.Filter(props.Search)

	for i, _ := range props.Selected {
		items[i].Selected = true
	}

	items[props.Cursor].Current = true

	c.items = items[props.Start:props.End]

	return nil
}

func (c *List) Render(w, h int) string {
	var items []string
	for _, i := range c.items {
		items = append(items, i.Render(w, h))
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}
