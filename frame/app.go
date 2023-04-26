package frame

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/confirm"
	"github.com/ohzqq/teacozy/frame/header"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
	"github.com/ohzqq/teacozy/util"
	"golang.org/x/exp/maps"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	router       *Router
	Routes       map[string]router.RouteInitializer
	defaultRoute string

	Confirm        confirm.Props
	confirmChoices bool

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
	editedVal   string
	choices     teacozy.Items
	keyMap      keys.KeyMap
	Style       Style

	Header *header.Component

	teacozy.Props
}

type Style struct {
	Confirm lipgloss.Style
	Footer  lipgloss.Style
	Header  lipgloss.Style
	Status  lipgloss.Style
}

func New(opts ...Option) *App {
	a := &App{
		router:   NewRouter(),
		Routes:   make(map[string]router.RouteInitializer),
		selected: make(map[int]struct{}),
		Style:    DefaultStyle(),
		cursor:   0,
		width:    util.TermWidth(),
		height:   util.TermHeight() - 2,
		limit:    10,
		Props:    teacozy.NewProps(),
	}

	//a.AddKey(keys.New("a").Cmd(keys.UpdateItem(keys.ToggleItems)))
	a.router.UpdateRoutes = a.UpdateRoutes

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (c *App) ItemProps() teacozy.Props {
	props := teacozy.NewProps()
	props.Paginator = c.Paginator
	props.Items = c.choices
	props.Selected = c.Selected
	props.SetKeyMap(c.Paginator.KeyMap)
	props.SetCurrent = c.SetCurrent
	props.ReadOnly = c.ReadOnly
	return props
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	switch len(c.Routes) {
	case 0:
		c.NewRoute(c)
	default:
		if c.defaultRoute != "" {
			c.Routes["default"] = c.Routes[c.defaultRoute]
		} else {
			k := maps.Keys(c.Routes)[0]
			c.Routes["default"] = c.Routes[k]
		}
	}
	c.Routes["confirm"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := confirm.New()
		p := c.Confirm
		return component, component.Init(p)
	}

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
	cmds = append(cmds, c.router.Init(c.Routes))
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
	//case keys.SaveChangesMsg:
	//fmt.Println("save edit")
	//return msg.Cmd(c.Current())
	//case keys.SaveEditMsg:
	case confirm.GetConfirmationMsg:
		//switch reactea.CurrentRoute() {
		//case "list":
		if !c.confirmChoices {
			return reactea.Destroy
		}
		//fallthrough
		//default:
		c.Confirm = msg.Props
		cmds = append(cmds, keys.ChangeRoute("confirm"))
		//}

	case keys.SaveChangesMsg:
		return keys.UpdateStatus("save main")
	//fmt.Println("save edit main")

	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height - 2
		return nil
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
	//var header string
	//if c.title != "" {
	//header = c.Style.Header.Render(c.title)
	//}
	return c.Header.Render(w, h)
}

func (c App) renderFooter(w, h int) string {
	var footer string

	footer = fmt.Sprintf(
		"cur route %v, per %v",
		reactea.CurrentRoute(),
		c.router.PrevRoute,
	)

	if c.footer != "" {
		footer = c.Style.Header.Render(c.footer)
	}

	return footer
}

func (c *App) NewRoute(r Route) {
	c.Routes[c.Name()] = r.Initializer(c.ItemProps())
}

func (c *App) UpdateRoutes(r Route) {
	c.NewRoute(r)
	c.router.Init(c.Routes)
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

func (c *App) Initialize(a *App) {
	a.Routes["default"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := reactea.Componentify[teacozy.Props](teacozy.Renderer)
		props := a.ItemProps()
		return component, component.Init(props)
	}
}

func (c *App) Initializer(props teacozy.Props) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := reactea.Componentify[teacozy.Props](teacozy.Renderer)
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

func DefaultStyle() Style {
	return Style{
		Confirm: lipgloss.NewStyle().Background(color.Red()).Foreground(color.Black()),
		Footer:  lipgloss.NewStyle().Foreground(color.Green()),
		Header:  lipgloss.NewStyle().Background(color.Purple()).Foreground(color.Black()),
		Status:  lipgloss.NewStyle().Foreground(color.Green()),
	}
}
