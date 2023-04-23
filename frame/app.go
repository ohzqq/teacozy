package frame

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
	"github.com/ohzqq/teacozy/util"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter   *Router
	Routes       map[string]router.RouteInitializer
	defaultRoute string

	start       int
	end         int
	width       int
	height      int
	selected    map[int]struct{}
	numSelected int
	limit       int
	CurrentItem int
	noLimit     bool
	cursor      int
	title       string
	header      string
	footer      string
	choices     teacozy.Items
	paginator   *pagy.Paginator
	Style       Style
}

type Style struct {
	Confirm lipgloss.Style
	Footer  lipgloss.Style
	Header  lipgloss.Style
	Status  lipgloss.Style
}

func New(opts ...Option) *App {
	a := &App{
		mainRouter: NewRouter(),
		Routes:     make(map[string]router.RouteInitializer),
		selected:   make(map[int]struct{}),
		Style:      DefaultStyle(),
		cursor:     0,
		width:      util.TermWidth(),
		height:     util.TermHeight() - 2,
		limit:      10,
	}
	a.mainRouter.UpdateRoutes = a.UpdateRoutes

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (c *App) ItemProps() teacozy.Props {
	props := teacozy.NewProps()
	props.Paginator = c.paginator
	props.Items = c.choices
	props.Selected = c.selected
	props.SetKeyMap(c.paginator.KeyMap)
	props.SetCurrent = c.SetCurrent
	return props
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	c.NewRoute(c)
	if c.defaultRoute != "" {
		c.Routes["default"] = c.Routes[c.defaultRoute]
	}

	if c.noLimit {
		c.limit = c.choices.Len()
	}

	c.paginator = pagy.New(c.height, c.choices.Len())
	c.paginator.SetKeyMap(DefaultKeyMap())

	var cmds []tea.Cmd
	cmds = append(cmds, keys.ChangeRoute("default"))
	cmds = append(cmds, c.mainRouter.Init(c.Routes))
	return tea.Batch(cmds...)
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case keys.ToggleItemMsg:
		c.ToggleItems(c.Current())
		cmds = append(cmds, keys.LineDown)
	case keys.StartEditingMsg:
		return keys.EditItem(c.Current())
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height - 2
		return nil
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	cmd = c.mainRouter.Update(msg)
	cmds = append(cmds, cmd)

	c.paginator, cmd = c.paginator.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *App) Render(w, h int) string {
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

	body := c.mainRouter.Render(c.width, height)
	view = append(view, body)

	if footer != "" {
		view = append(view, footer)
	}

	return lipgloss.JoinVertical(lipgloss.Left, view...)
}

func (c App) renderHeader(w, h int) string {
	var header string
	if c.title != "" {
		header = c.Style.Header.Render(c.title)
	}
	return header
}

func (c App) renderFooter(w, h int) string {
	var footer string

	footer = fmt.Sprintf(
		"cur route %v, per %v, current %v",
		reactea.CurrentRoute(),
		c.paginator.Current(),
		c.Current(),
	)

	if c.footer != "" {
		footer = c.Style.Header.Render(c.footer)
	}

	return footer
}

func (c *App) NewRoute(r Route) {
	r.Initialize(c)
}

func (c *App) UpdateRoutes(r Route) {
	c.NewRoute(r)
	c.mainRouter.Init(c.Routes)
}

func (c *App) SetKeyMap(km keys.KeyMap) *App {
	c.paginator.SetKeyMap(km)
	return c
}

func (m *App) ToggleItems(items ...int) {
	for _, idx := range items {
		m.CurrentItem = idx
		if _, ok := m.selected[idx]; ok {
			delete(m.selected, idx)
			m.numSelected--
		} else if m.numSelected < m.limit {
			m.selected[idx] = struct{}{}
			m.numSelected++
		}
	}
}

func (m App) Chosen() []map[string]string {
	var chosen []map[string]string
	if len(m.selected) > 0 {
		for k := range m.selected {
			l := m.choices.Label(k)
			v := m.choices.String(k)
			chosen = append(chosen, map[string]string{l: v})
		}
	}
	return chosen
}

func (m *App) SetCurrent(idx int) {
	m.CurrentItem = idx
	m.paginator.SetCurrent(idx)
}

func (m *App) Current() int {
	return m.CurrentItem
}

func (c App) Selected() {
	fmt.Printf("sel %+V\n", c.selected)
}

func (c *App) Initialize(a *App) {
	a.Routes["default"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := reactea.Componentify[teacozy.Props](teacozy.Renderer)
		return component, component.Init(a.ItemProps())
	}
}

func (c App) Name() string {
	return "default"
}

func (c *App) SetWidth(n int) *App {
	c.width = n
	return c
}

func (c *App) SetHeight(n int) *App {
	c.height = n
	return c
}

func (c *App) SetSize(w, h int) *App {
	c.width = w
	c.height = h
	return c
}

func DefaultKeyMap() keys.KeyMap {
	return keys.KeyMap{
		keys.Up().AddKeys("k"),
		keys.Down().AddKeys("j"),
		keys.HalfPgUp().AddKeys("K"),
		keys.HalfPgDown().AddKeys("J"),
		keys.Home().AddKeys("g"),
		keys.End().AddKeys("G"),
		keys.Quit().AddKeys("q"),
	}
}

func DefaultStyle() Style {
	return Style{
		Confirm: lipgloss.NewStyle().Background(color.Red()).Foreground(color.Black()),
		Footer:  lipgloss.NewStyle().Foreground(color.Green()),
		Header:  lipgloss.NewStyle().Background(color.Purple()).Foreground(color.Black()),
		Status:  lipgloss.NewStyle().Foreground(color.Green()),
	}
}
