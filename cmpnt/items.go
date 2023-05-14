package cmpnt

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy"
	"github.com/sahilm/fuzzy"
)

type Items struct {
	ReadOnly    bool
	SetCurrent  func(int)
	Current     func() int
	Highlighted func() int
	IsSelected  func(int) bool
	Style       Style
	Items       teacozy.Items
	Matches     fuzzy.Matches
}

func NewItems(items teacozy.Items) Items {
	p := Items{
		Items: items,
		Style: DefaultStyle(),
	}
	return p
}

func (props Items) Copy() Items {
	items := NewItems(props.Items)
	items.IsSelected = props.IsSelected
	return items
}

func (props Items) Render() string {
	var s strings.Builder
	for i, m := range props.Matches {
		cur := i == props.Highlighted()
		sel := props.IsSelected(m.Index)
		label := props.Items.Label(m.Index)

		pre := props.PrefixText(label, sel, cur)
		style := props.PrefixStyle(label, sel, cur)

		// only print the prefix if it's a list or there's a label
		if !props.ReadOnly || label != "" {
			s.WriteString(style.Render(pre))
		}

		// render the rest of the line
		text := lipgloss.StyleRunes(
			m.Str,
			m.MatchedIndexes,
			props.Style.Match,
			props.Style.Normal.Style,
		)

		s.WriteString(lipgloss.NewStyle().Render(text))
		s.WriteString("\n")
	}

	return s.String()
}

func (c *Items) SetMatches(matches fuzzy.Matches) {
	c.Matches = matches
}

func (c Items) PrefixText(label string, selected, current bool) string {
	switch {
	case label != "":
		return label
	case current:
		return c.Style.Cursor.Text
	case selected && !c.ReadOnly:
		return c.Style.Selected.Text
	default:
		return c.Style.Normal.Text
	}
}

func (c Items) PrefixStyle(label string, selected, current bool) Prefix {
	switch {
	case current:
		return c.Style.Cursor
	case selected && !c.ReadOnly:
		return c.Style.Selected
	case label != "":
		return c.Style.Label
	default:
		return c.Style.Normal
	}
}

func (c Items) Len() int {
	return c.Items.Len()
}

func (c Items) String(idx int) string {
	return c.Items.String(idx)
}

func (c Items) Label(idx int) string {
	return c.Items.Label(idx)
}
