package teacozy

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
	"github.com/sahilm/fuzzy"
)

type State struct {
	*pagy.Paginator
	name       string
	Items      Items
	Selected   map[int]struct{}
	Search     string
	ReadOnly   bool
	SetCurrent func(int)
	SetHelp    func(keys.KeyMap)
	Style      Style
}

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

type Items interface {
	fuzzy.Source
	Label(int) string
	Set(int, string)
}

func NewProps(items Items) State {
	p := State{
		Items:    items,
		Selected: make(map[int]struct{}),
		Style:    DefaultStyle(),
	}
	return p
}

func Renderer(props State, w, h int) string {
	var s strings.Builder
	h = h - 2

	// get matched items
	items := props.ExactMatches(props.Search)

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

func (p State) Initializer(props State) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := reactea.Componentify[State](Renderer)
		return component, component.Init(props)
	}
}

func (p State) Name() string {
	return p.name
}

func (p *State) SetName(name string) {
	p.name = name
}

func (c State) PrefixText(label string, selected, current bool) string {
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

func (c State) PrefixStyle(label string, selected, current bool) Prefix {
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

func (c *State) Filter(s string, w, h int) string {
	c.Search = s
	c.ResetCursor()
	p := *c
	return Renderer(p, w, h)
}

func (c *State) ExactMatches(search string) fuzzy.Matches {
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
