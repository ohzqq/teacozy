package item

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
)

type List struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]
}

type Props struct {
	Choices  Choices
	Selected map[int]struct{}
	Cursor   int
	Start    int
	End      int
	Search   string
}

func NewList() *List {
	return &List{}
}

func Renderer(props Props, w, h int) string {
	items := props.Choices.Filter(props.Search)

	for i, _ := range props.Selected {
		items[i].Selected = true
	}

	items[props.Cursor].Current = true

	var rendered []string
	for _, i := range items[props.Start:props.End] {
		rendered = append(rendered, i.Render(w, h))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}

func (p *Props) SetCursor(n int) {
	p.Cursor = n
}

func (c *List) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	return nil
}

//func (c List) Update(msg tea.Msg) tea.Cmd {
//}

func (c *List) Render(w, h int) string {
	items := c.Props().Choices.Filter(c.Props().Search)

	for i, _ := range c.Props().Selected {
		items[i].Selected = true
	}

	items[c.Props().Cursor].Current = true

	var rendered []string
	for _, i := range items[c.Props().Start:c.Props().End] {
		rendered = append(rendered, i.Render(w, h))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}

func (c Props) Matches() []Item {
	return c.Choices.Filter(c.Search)
}
