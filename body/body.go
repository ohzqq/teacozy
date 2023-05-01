package body

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
	"github.com/ohzqq/teacozy/util"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[teacozy.Props]

	width  int
	height int

	CurrentItem int

	keyMap keys.KeyMap

	*pagy.Paginator
}

type Props struct {
	teacozy.Props
	Width  int
	Height int
}

func New() *Component {
	c := &Component{
		width:  util.TermWidth(),
		height: util.TermHeight() - 2,
	}

	return c
}

func (c *Component) Init(props teacozy.Props) tea.Cmd {
	c.UpdateProps(props)

	c.Paginator = pagy.New(c.height, c.Props().Items.Len())
	c.Paginator.SetKeyMap(keys.VimKeyMap())

	return nil
}

func (c *Component) Initializer(props teacozy.Props) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		return component, component.Init(props)
	}
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}

		for _, k := range c.keyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	p := c.Props().Paginator

	p, cmd = p.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c Component) Name() string {
	return "body"
}

func (c *Component) Render(w, h int) string {
	var s strings.Builder
	h = h - 2

	// get matched items
	p := c.Props()
	items := p.ExactMatches(c.Props().Search)

	c.Props().SetPerPage(h)

	// update the total items based on the filter results, this prevents from
	// going out of bounds
	c.Props().SetTotal(len(items))

	for i, m := range items[c.Props().Start():c.Props().End()] {
		var cur bool
		if i == c.Props().Highlighted() {
			c.Props().SetCurrent(m.Index)
			cur = true
		}

		var sel bool
		if _, ok := c.Props().Selected[m.Index]; ok {
			sel = true
		}

		label := c.Props().Items.Label(m.Index)
		pre := c.Props().PrefixText(label, sel, cur)
		style := c.Props().PrefixStyle(label, sel, cur)

		// only print the prefix if it's a list or there's a label
		if !c.Props().ReadOnly || label != "" {
			s.WriteString(style.Render(pre))
		}

		// render the rest of the line
		text := lipgloss.StyleRunes(
			m.Str,
			m.MatchedIndexes,
			c.Props().Style.Match,
			c.Props().Style.Normal.Style,
		)

		s.WriteString(lipgloss.NewStyle().Render(text))
		s.WriteString("\n")
	}

	return s.String()
}

func (c *Component) SetKeyMap(km keys.KeyMap) *Component {
	c.Paginator.SetKeyMap(km)
	return c
}

func (c *Component) AddKey(k *keys.Binding) *Component {
	if !c.keyMap.Contains(k) {
		c.keyMap.AddBinds(k)
	} else {
		c.keyMap.Replace(k)
	}
	return c
}

func (m *Component) SetCurrent(idx int) {
	m.CurrentItem = idx
}

func (m *Component) Current() int {
	return m.CurrentItem
}
