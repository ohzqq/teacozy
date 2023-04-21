package teacozy

import (
	"strings"

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
	Prefix   Prefixes
}

type Prefix struct {
	Text  string
	Style lipgloss.Style
}

type Prefixes struct {
	Cursor     Prefix
	Selected   Prefix
	Unselected Prefix
}

type Items interface {
	Find(string) fuzzy.Matches
}

func NewProps() Props {
	d := item.DefaultStyle()
	return Props{
		Selected: make(map[int]struct{}),
		Prefix: Prefixes{
			Cursor: Prefix{
				Text:  "x",
				Style: d.Cursor,
			},
			Selected: Prefix{
				Text:  "x",
				Style: d.Selected,
			},
			Unselected: Prefix{
				Text:  " ",
				Style: d.Unselected,
			},
		},
	}
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
		var s strings.Builder
		//rendered = append(rendered, i.Render(w, h))
		pre := "x"

		if i.Label != "" {
			pre = i.Label
		}

		switch {
		case i.Current:
			pre = i.Style.Cursor.Render(pre)
		default:
			if i.Selected {
				pre = i.Style.Selected.Render(pre)
			} else if i.Label == "" {
				pre = strings.Repeat(" ", lipgloss.Width(pre))
			} else {
				pre = i.Style.Label.Render(pre)
			}
		}

		s.WriteString("[")
		s.WriteString(pre)
		s.WriteString("]")

		text := lipgloss.StyleRunes(
			i.Str,
			i.MatchedIndexes,
			i.Style.Match,
			i.Style.Unselected,
		)
		s.WriteString(lipgloss.NewStyle().Render(text))

		rendered = append(rendered, s.String())
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
