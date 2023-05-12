package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/header"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
	"github.com/ohzqq/teacozy/util"
)

type Page struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	confirmChoices bool
	readOnly       bool

	width  int
	height int

	numSelected int
	limit       int
	CurrentItem int
	noLimit     bool

	footer string

	choices teacozy.Items
	keyMap  keys.KeyMap
	Style   AppStyle

	Header *header.Component
	title  string
	header string

	help keys.KeyMap

	teacozy.State
}

type AppStyle struct {
	Footer lipgloss.Style
}

func New(opts ...Option) *Page {
	c := &Page{
		width:  util.TermWidth(),
		height: util.TermHeight() - 2,
		limit:  10,
	}

	c.Style = AppStyle{
		Footer: lipgloss.NewStyle().Foreground(color.Green()),
	}

	for _, opt := range opts {
		opt(c)
	}

	c.State = teacozy.NewProps(c.choices)
	c.State.SetCurrent = c.SetCurrent
	c.State.SetHelp = c.SetHelp
	c.State.ReadOnly = c.readOnly

	return c
}

func (c *Page) Init(reactea.NoProps) tea.Cmd {
	if c.noLimit {
		c.limit = c.choices.Len()
	}

	if !c.ReadOnly {
		c.AddKey(keys.Toggle().AddKeys(" "))
	}

	c.Paginator = pagy.New(c.height, c.choices.Len())
	c.Paginator.SetKeyMap(keys.VimKeyMap())

	c.Header = header.New()
	c.Header.Init(
		header.Props{
			Title: c.title,
		},
	)

	c.SetHelp(c.Paginator.KeyMap)
	c.AddKey(keys.Help())

	return nil
}

func (c *Page) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case keys.ShowHelpMsg:
		cmds = append(cmds, keys.ChangeRoute("help"))
		//help := NewProps(c.help)
		//help.SetName("help")
		//return ChangeRoute(&help)

	case keys.UpdateItemMsg:
		return msg.Cmd(c.Current())

	case keys.ToggleItemsMsg, keys.ToggleItemMsg:
		c.ToggleItems(c.Current())
		cmds = append(cmds, keys.LineDown)

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

	c.Paginator, cmd = c.Paginator.Update(msg)
	cmds = append(cmds, cmd)

	cmd = c.Header.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *Page) Render(w, h int) string {
	height := c.height
	var view []string

	if head := c.renderHeader(w, height); head != "" {
		height -= lipgloss.Height(head)
		view = append(view, head)
	}

	footer := c.renderFooter(w, height)
	if footer != "" {
		height -= lipgloss.Height(footer)
	}

	body := teacozy.Renderer(c.State, c.width, height)
	view = append(view, body)

	if footer != "" {
		view = append(view, footer)
	}

	return lipgloss.JoinVertical(lipgloss.Left, view...)
}

func (c Page) renderHeader(w, h int) string {
	return c.Header.Render(w, h)
}

func (c Page) renderFooter(w, h int) string {
	var footer string

	//footer = fmt.Sprintf(
	//"cur route %v, per %v",
	//reactea.CurrentRoute(),
	//c.router.PrevRoute,
	//)

	if c.footer != "" {
		footer = c.Style.Footer.Render(c.footer)
	}

	return footer
}

func (c *Page) SetKeyMap(km keys.KeyMap) *Page {
	c.Paginator.SetKeyMap(km)
	return c
}

func (c Page) ToggleItem() {
	c.ToggleItems(c.Current())
}

func (c *Page) AddKey(k *keys.Binding) *Page {
	if !c.keyMap.Contains(k) {
		c.keyMap.AddBinds(k)
	} else {
		c.keyMap.Replace(k)
	}
	return c
}

func (c *Page) ToggleItems(items ...int) {
	for _, idx := range items {
		c.CurrentItem = idx
		if _, ok := c.Selected[idx]; ok {
			delete(c.Selected, idx)
			c.numSelected--
		} else if c.numSelected < c.limit {
			c.Selected[idx] = struct{}{}
			c.numSelected++
		}
	}
}

func (m Page) Chosen() []map[string]string {
	var chosen []map[string]string
	if len(m.Selected) > 0 {
		for k := range m.Selected {
			l := m.choices.Label(k)
			v := m.choices.String(k)
			chosen = append(chosen, map[string]string{l: v})
		}
	}
	return chosen
}

func (m *Page) SetCurrent(idx int) {
	m.CurrentItem = idx
}

func (m *Page) Current() int {
	return m.CurrentItem
}

func (c *Page) SetWidth(n int) *Page {
	c.width = n
	return c
}

func (c *Page) SetHelp(km keys.KeyMap) {
	c.help = km
}

func (c *Page) SetHeight(n int) *Page {
	c.height = n
	return c
}

func (c *Page) SetSize(w, h int) *Page {
	c.width = w
	c.height = h
	return c
}

func (c *Page) SetHeader(h string) {
	c.header = h
}
