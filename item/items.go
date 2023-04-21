package item

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/pagy"
)

type List struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]
}

type Props struct {
	*pagy.Paginator
	Choices  Choices
	Selected map[int]struct{}
	Search   string
}

func NewList() *List {
	return &List{}
}

func Renderer(props Props, w, h int) string {
	items := props.Choices.Filter(props.Search)
	props.SetTotal(len(items))

	for i, _ := range props.Selected {
		items[i].Selected = true
	}

	items[props.Cursor()].Current = true

	var rendered []string
	for _, i := range items[props.Start():props.End()] {
		rendered = append(rendered, i.Render(w, h))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}

func (c *Props) Filter(s string) {
	c.Search = s
}

func (c Props) Matches() []Item {
	return c.Choices.Filter(c.Search)
}

func (c *Props) SetPaginator(p *pagy.Paginator) {
	c.Paginator = p
}
