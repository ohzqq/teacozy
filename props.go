package teacozy

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/pagy"
	"github.com/sahilm/fuzzy"
)

type Props struct {
	*pagy.Paginator
	Choices  item.Choices
	Items    Items
	Selected map[int]struct{}
	Search   string
	Prefix   struct {
		Cursor     Prefix
		Selected   Prefix
		Unselected Prefix
	}
}

type Prefix struct {
	Text  string
	Style lipgloss.Style
}

type Items interface {
	Find(string) fuzzy.Matches
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

func (p Prefix) Render(t ...string) string {
	text := p.Text
	if len(t) > 0 {
		text = t[0]
	}
	return p.Style.Render(text)
}
