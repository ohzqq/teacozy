package cmpnt

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
)

type Prefix struct {
	Fmt   string
	Text  string
	Style lipgloss.Style
}

type Style struct {
	Cursor   Prefix
	Label    Prefix
	Normal   Prefix
	Selected Prefix
	Match    lipgloss.Style
}

func (c Pager) PrefixText(label string, selected, current bool) string {
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

func (c Pager) PrefixStyle(label string, selected, current bool) Prefix {
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

func (p Prefix) Render(pre ...string) string {
	text := p.Text
	if len(pre) > 0 {
		if t := pre[0]; t != "" {
			text = t
		}
	}
	return fmt.Sprintf(p.Fmt, p.Style.Render(text))
}

func DefaultStyle() Style {
	return Style{
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
