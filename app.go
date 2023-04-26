package teacozy

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/header"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
	"github.com/ohzqq/teacozy/util"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	router       *Router
	Routes       map[string]Route
	defaultRoute string

	confirmChoices bool

	width  int
	height int

	numSelected int
	limit       int
	CurrentItem int
	noLimit     bool

	footer string

	choices Items
	keyMap  keys.KeyMap
	Style   AppStyle

	Header *header.Component
	title  string
	header string

	Props
}

type AppStyle struct {
	Footer lipgloss.Style
}

func New(opts ...Option) *App {
	a := &App{
		router:       NewRouter(),
		Routes:       make(map[string]Route),
		defaultRoute: "default",
		Props:        NewProps(),
		width:        util.TermWidth(),
		height:       util.TermHeight() - 2,
		limit:        10,
	}

	a.Style = AppStyle{
		Footer: lipgloss.NewStyle().Foreground(color.Green()),
	}

	a.router.UpdateRoutes = a.UpdateRoute

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (c *App) ItemProps() Props {
	props := NewProps()
	props.Paginator = c.Paginator
	props.Items = c.choices
	props.Selected = c.Selected
	props.SetKeyMap(c.Paginator.KeyMap)
	props.SetCurrent = c.SetCurrent
	props.ReadOnly = c.ReadOnly
	return props
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	c.Routes["default"] = c

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

	var cmds []tea.Cmd
	cmds = append(cmds, keys.ChangeRoute("default"))
	cmds = append(cmds, c.InitRoutes())
	return tea.Batch(cmds...)
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
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

	cmd = c.router.Update(msg)
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

	body := c.router.Render(c.width, height)
	view = append(view, body)

	if footer != "" {
		view = append(view, footer)
	}

	return lipgloss.JoinVertical(lipgloss.Left, view...)
}

func (c App) renderHeader(w, h int) string {
	return c.Header.Render(w, h)
}

func (c App) renderFooter(w, h int) string {
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

func (c *App) InitRoutes() tea.Cmd {
	routes := make(map[string]router.RouteInitializer, len(c.Routes))
	for name, route := range c.Routes {
		routes[name] = route.Initializer(c.ItemProps())
	}

	p := RouterProps{
		Routes:      routes,
		Default:     c.defaultRoute,
		ChangeRoute: c.UpdateRoute,
	}

	return c.router.Init(p)
}

func (c *App) UpdateRoute(r Route) tea.Cmd {
	c.Routes[c.Name()] = r
	c.InitRoutes()
	return keys.ChangeRoute(r.Name())
}

func (c *App) SetKeyMap(km keys.KeyMap) *App {
	c.Paginator.SetKeyMap(km)
	return c
}

func (c App) ToggleItem() {
	c.ToggleItems(c.Current())
}

func (c *App) AddKey(k *keys.Binding) *App {
	if !c.keyMap.Contains(k) {
		c.keyMap.AddBinds(k)
	} else {
		c.keyMap.Replace(k)
	}
	return c
}

func (c *App) ToggleItems(items ...int) {
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

func (m App) Chosen() []map[string]string {
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

func (m *App) SetCurrent(idx int) {
	m.CurrentItem = idx
}

func (m *App) Current() int {
	return m.CurrentItem
}

func (c *App) Initializer(props Props) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := reactea.Componentify[Props](Renderer)
		return component, component.Init(props)
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

func (c *App) SetHeader(h string) {
	c.header = h
}
