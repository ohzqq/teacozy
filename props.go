package teacozy

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
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
	SetCurrent func(int)
	Style      PropsStyle
}

type Prefix struct {
	Fmt   string
	Text  string
	Style lipgloss.Style
}

type PropsStyle struct {
	Cursor   Prefix
	Label    Prefix
	Normal   Prefix
	Selected Prefix
	Match    lipgloss.Style
}

type Items interface {
	fuzzy.Source
	Label(int) string
	Set(int, string)
}

func NewProps() Props {
	p := Props{
		Selected: make(map[int]struct{}),
		Style:    DefaultPropsStyle(),
	}
	return p
}

func Renderer(props Props, w, h int) string {
	var s strings.Builder
	h = h - 2

	// get matched items
	items := props.exactMatches(props.Search)

	props.SetPerPage(h)

	// update the total items based on the filter results, this prevents from
	// going out of bounds
	props.SetTotal(len(items))

	for i, m := range items[props.Start():props.End()] {
		var cur bool
		if i == props.Highlighted() {
			props.SetCurrent(m.Index)
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
			props.Style.Normal.Style,
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
		return c.Style.Cursor.Text
	case selected && c.Selectable:
		return c.Style.Selected.Text
	default:
		return c.Style.Normal.Text
	}
}

func (c Props) prefixStyle(label string, selected, current bool) Prefix {
	switch {
	case current:
		return c.Style.Cursor
	case selected && c.Selectable:
		return c.Style.Selected
	case label != "":
		return c.Style.Label
	default:
		return c.Style.Normal
	}
}

func (c *Props) Filter(s string) {
	c.Search = s
	c.ResetCursor()
}

func (c *Props) exactMatches(search string) fuzzy.Matches {
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

func DefaultPropsStyle() PropsStyle {
	return PropsStyle{
		Match: lipgloss.NewStyle().Foreground(color.Cyan()),
		Cursor: Prefix{
			Fmt:   currentFmt,
			Text:  currentPrefix,
			Style: lipgloss.NewStyle().Foreground(color.Green()),
		},
		Selected: Prefix{
			Fmt:   selectedFmt,
			Text:  selectedPrefix,
			Style: lipgloss.NewStyle().Foreground(color.Grey()),
		},
		Normal: Prefix{
			Fmt:   unselectedFmt,
			Text:  unselectedPrefix,
			Style: lipgloss.NewStyle().Foreground(color.Fg()),
		},
		Label: Prefix{
			Fmt:   labelFmt,
			Style: lipgloss.NewStyle().Foreground(color.Purple()),
		},
	}
}

const (
	selectedPrefix   = "x"
	selectedFmt      = "[%s]"
	unselectedPrefix = " "
	unselectedFmt    = "[%s]"
	currentPrefix    = "x"
	currentFmt       = "[%s]"
	labelFmt         = "[%s]"
)
