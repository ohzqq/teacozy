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
	Choices    item.Choices
	Items      Source
	Selected   map[int]struct{}
	Search     string
	Selectable bool
	Prefix     Prefixes
	Style      item.Style
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

type Source interface {
	fuzzy.Source
	Label(int) string
	Set(int, string)
}

func NewProps() Props {
	d := item.DefaultStyle()
	return Props{
		Style:    d,
		Selected: make(map[int]struct{}),
		Prefix: Prefixes{
			Cursor: Prefix{
				Text:  "[x]",
				Style: d.Cursor,
			},
			Selected: Prefix{
				Text:  "[x]",
				Style: d.Selected,
			},
			Unselected: Prefix{
				Text:  "[ ]",
				Style: d.Unselected,
			},
		},
	}
}

func Renderer(props Props, w, h int) string {
	items := props.exactMatches(props.Search)
	props.SetTotal(len(items))

	var rendered []string
	for _, m := range items[props.Start():props.End()] {
		var s strings.Builder

		var pre string
		label := props.Items.Label(m.Index)
		if label != "" {
			pre = label
		}
		switch {
		case m.Index == props.Paginator.Cursor():
			pre = props.Prefix.Cursor.Render(pre)
		default:
			if props.Selectable {
				if _, ok := props.Selected[m.Index]; ok {
					pre = props.Prefix.Selected.Render(pre)
				} else {
					pre = props.Prefix.Unselected.Render(pre)
				}
			}
			if label != "" {
				pre = props.Style.Label.Render(pre)
			}
		}
		if props.Selectable || label != "" {
			s.WriteString(pre)
		}

		text := lipgloss.StyleRunes(
			m.Str,
			m.MatchedIndexes,
			props.Style.Match,
			props.Style.Unselected,
		)
		s.WriteString(lipgloss.NewStyle().Render(text))

		rendered = append(rendered, s.String())
	}

	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}

func (c *Props) Filter(s string) {
	c.Search = s
}

func (c Props) exactMatches(search string) fuzzy.Matches {
	if search != "" {
		if m := fuzzy.FindFrom(search, c.Items); len(m) > 0 {
			return m
		}
	}
	return SourceToMatches(c.Items)
}

func (p Prefix) Render(pre ...string) string {
	text := p.Text
	if len(pre) > 0 {
		if t := pre[0]; t != "" {
			text = t
		}
	}
	return p.Style.Render(text)
}

func SourceToMatches(src Source) fuzzy.Matches {
	items := make(fuzzy.Matches, src.Len())

	for i := 0; i < src.Len(); i++ {
		m := fuzzy.Match{
			Str:   src.String(i),
			Index: i,
		}
		items[i] = m
	}
	return items
}
