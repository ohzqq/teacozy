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
	Items    Source
	Selected map[int]struct{}
	Search   string
	Prefix   Prefixes
	Style    item.Style
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

type Items struct {
	src   Source
	items []Item
}

type Item struct {
	fuzzy.Match
	Label string
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

func OldRenderer(props Props, w, h int) string {
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

func Renderer(props Props, w, h int) string {
	items := props.exactMatches(props.Search)
	props.SetTotal(items.Len())

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
			//pre = props.Style.Cursor.Render(pre)
			pre = props.Prefix.Cursor.Render(pre)
		default:
			if _, ok := props.Selected[m.Index]; ok {
				//pre = props.Style.Selected.Render(pre)
				pre = props.Prefix.Selected.Render(pre)
			} else {
				pre = props.Prefix.Unselected.Render(pre)
				//pre = strings.Repeat(" ", lipgloss.Width(props.Prefix.Cursor.Text))
				//} else if i.Label == "" {
				//  pre = strings.Repeat(" ", lipgloss.Width(pre))
				//} else {
				//  pre = i.Style.Label.Render(pre)
			}
		}

		//s.WriteString("[")
		s.WriteString(pre)
		//s.WriteString("]")

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
		return fuzzy.FindFrom(search, c.Items)
	}
	return SourceToMatches(c.Items)
}

func (p Prefix) Render(t ...string) string {
	text := p.Text
	if len(t) > 0 {
		text = t[0]
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

func (i Items) String(idx int) string {
	return i.items[idx].Str
}

func (i Items) Label(idx int) string {
	return i.items[idx].Label
}

func (i Items) Len() int {
	return len(i.items)
}

func (i *Items) Set(idx int, val string) {
	i.src.Set(idx, val)
}
