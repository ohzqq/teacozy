package cmpnt

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
	"github.com/sahilm/fuzzy"
)

type Items struct {
	name        string
	Items       teacozy.Items
	Selected    map[int]struct{}
	Search      string
	ReadOnly    bool
	SetCurrent  func(int)
	SetPerPage  func(int)
	SetTotal    func(int)
	SetHelp     func(keys.KeyMap)
	Start       int
	End         int
	Highlighted int
	Style       Style
}

type ItemProps struct {
	ReadOnly    bool
	Highlighted int
	Style       Style
	Selected    map[int]struct{}
	Items       teacozy.Items
	Matches     fuzzy.Matches
	SetCurrent  func(int)
	Current     func() int
}

func NewItems(items teacozy.Items) ItemProps {
	p := ItemProps{
		Items:    items,
		Style:    DefaultStyle(),
		Selected: make(map[int]struct{}),
	}
	return p
}

func (props ItemProps) Render() string {
	var s strings.Builder
	for i, m := range props.Matches {
		var cur bool
		if i == props.Highlighted {
			props.SetCurrent(m.Index)
			cur = true
		}

		var sel bool
		if _, ok := props.Selected[m.Index]; ok {
			sel = true
		}

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

func (c ItemProps) PrefixText(label string, selected, current bool) string {
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

func (c ItemProps) PrefixStyle(label string, selected, current bool) Prefix {
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

func (c *Items) ExactMatches(search string) fuzzy.Matches {
	if search != "" {
		if m := fuzzy.FindFrom(search, c.Items); len(m) > 0 {
			return m
		}
	}
	return teacozy.SourceToMatches(c.Items)
}
