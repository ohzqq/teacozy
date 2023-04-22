package teacozy

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/pagy"
	"github.com/sahilm/fuzzy"
)

type Props struct {
	*pagy.Paginator
	Choices    item.Choices
	Items      Items
	Selected   map[int]struct{}
	Search     string
	Selectable bool
	Prefix     Prefixes
	Style      item.Style
}

type Prefix struct {
	Fmt   string
	Text  string
	Style lipgloss.Style
}

type Prefixes struct {
	Cursor     Prefix
	Selected   Prefix
	Unselected Prefix
	Label      Prefix
}

type Items interface {
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
				Fmt:   "[%s]",
				Text:  "x",
				Style: d.Cursor,
			},
			Selected: Prefix{
				Fmt:   "[%s]",
				Text:  "x",
				Style: d.Selected,
			},
			Unselected: Prefix{
				Fmt:   "[%s]",
				Text:  " ",
				Style: d.Unselected,
			},
			Label: Prefix{
				Fmt:   "[%s]",
				Style: d.Label,
			},
		},
	}
}

func Renderer(props Props, w, h int) string {
	var s strings.Builder

	// get matched items
	items := props.exactMatches(props.Search)

	// update the total items based on the filter results, this prevents from
	// going out of bounds
	props.SetTotal(len(items))

	for _, m := range items[props.Start():props.End()] {
		var cur bool
		if m.Index == items[props.Cursor()].Index {
			cur = true
		}

		var sel bool
		if _, ok := props.Selected[m.Index]; ok {
			sel = true
		}

		label := props.Items.Label(m.Index)
		pre := props.prefixText(label, sel, cur)
		style := props.prefixStyle(label, sel, cur)

		// only print the prefix if it's a list or there's a label
		if props.Selectable || label != "" {
			s.WriteString(style.Render(pre))
		}

		// render the rest of the line
		text := lipgloss.StyleRunes(
			m.Str,
			m.MatchedIndexes,
			props.Style.Match,
			props.Style.Unselected,
		)

		s.WriteString(lipgloss.NewStyle().Render(text))
		s.WriteString("\n")
	}

	return s.String()
}

func (c Props) prefixText(label string, selected, current bool) string {
	switch {
	case label != "":
		return label
	case current:
		return c.Prefix.Cursor.Text
	case selected && c.Selectable:
		return c.Prefix.Selected.Text
	default:
		return c.Prefix.Unselected.Text
	}
}

func (c Props) prefixStyle(label string, selected, current bool) Prefix {
	switch {
	case current:
		return c.Prefix.Cursor
	case selected && c.Selectable:
		return c.Prefix.Selected
	case label != "":
		return c.Prefix.Label
	default:
		return c.Prefix.Unselected
	}
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
	return fmt.Sprintf(p.Fmt, p.Style.Render(text))
}

func SourceToMatches(src Items) fuzzy.Matches {
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
