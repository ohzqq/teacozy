package app

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/app/header"
	"github.com/ohzqq/teacozy/app/input"
	"github.com/ohzqq/teacozy/app/item"
	"github.com/ohzqq/teacozy/app/list"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	Style       style.App
	width       int
	height      int
	Choices     []map[string]string
	Selected    map[int]struct{}
	NumSelected int
	Limit       int
	search      string
	inText      string
	list        *list.Component
	input       *input.Component
	footer      string
	header      string
	status      string

	PrevRoute string
	routes    map[string]reactea.SomeComponent

	Viewport viewport.Model
}

func New(choices []string) *App {
	a := &App{
		mainRouter: router.New(),
		width:      util.TermWidth(),
		height:     10,
		Choices:    item.MapChoices(choices),
		Style:      style.DefaultAppStyle(),
		Selected:   make(map[int]struct{}),
		Limit:      1,
		routes:     make(map[string]reactea.SomeComponent),
		//header:     "poot",
	}

	return a
}

func (c *App) listProps() list.Props {
	p := list.Props{
		Matches:     Filter(c.search, c.Choices),
		Selected:    c.Selected,
		ToggleItems: c.ToggleItems,
	}
	return p
}

func (c App) Height() int {
	return c.height
}

func (c App) Width() int {
	return c.width
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	c.list = list.New()
	c.list.SetKeyMap(keys.VimListKeyMap())
	c.list.Init(c.listProps())
	//c.input = input.New()
	//c.input.Init(input.Props{
	//Filter: c.Input,
	//})

	//c.routes["list"] = c.list
	//c.routes["filter"] = c.input
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := new(struct {
				reactea.BasicComponent
				reactea.InvisibleComponent
			})
			return component, nil
		},
		"status": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			//component := reactea.Componentify[string](RenderHeader)
			component := header.New()
			//return component, nil
			return component, component.Init(header.Props{Msg: c.status})
		},
		"filter": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := input.New()
			c.list.SetKeyMap(keys.DefaultListKeyMap())
			return component, component.Init(input.Props{Filter: c.Input})
		},
	})
}

func (c *App) SetHeader(h string) {
	c.header = h
}

func (c *App) SetStatus(h string) {
	c.status = h
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	if reactea.CurrentRoute() == "" {
		reactea.SetCurrentRoute("list")
	}
	reactea.AfterUpdate(c)
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case header.StatusMsg:
		c.SetStatus(msg.Status)
		cmds = append(cmds, message.ChangeRoute("status"))

	case message.StopFilteringMsg:
		c.Input("")
		c.list.SetKeyMap(keys.VimListKeyMap())
		cmds = append(cmds, message.ChangeRoute("list"))

	case message.StartFilteringMsg:
		//c.list.SetKeyMap(keys.DefaultListKeyMap())
		cmds = append(cmds, message.ChangeRoute("filter"))

	case message.ConfirmMsg:
		//c.ConfirmAction = ""
	case message.GetConfirmationMsg:
		//c.ConfirmAction = msg.Question
	case message.ChangeRouteMsg:
		route := msg.Name
		switch route {
		case "list":
			c.list.SetCursor(0)
		case "filter":
			//c.Footer("")
		case "prev":
			route = c.PrevRoute
		case "help":
			//p := c.NewProps(KeymapToProps(c.help))
			//p.Height = c.Items.Height
			//p.Width = c.Items.Width
			//c.Routes["help"] = view.New().Initializer(p)
		}
		c.ChangeRoute(route)
	case message.ReturnSelectionsMsg:
		switch reactea.CurrentRoute() {
		case "filter":
		default:
			return reactea.Destroy
		}
	case message.QuitMsg:
		return reactea.Destroy
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+o":
			return header.StatusUpdate("toot")
		case "/":
			cmds = append(cmds, message.StartFiltering())
		case "ctrl+c":
			return reactea.Destroy
		case "y":
			cmds = append(cmds, message.Confirm(true))
		case "n":
			cmds = append(cmds, message.Confirm(false))
		}
	}
	switch reactea.CurrentRoute() {
	case "filter":
		c.search = c.inText
	case "help":
	}

	cmds = append(cmds, c.list.Update(msg))
	cmds = append(cmds, c.mainRouter.Update(msg))

	return tea.Batch(cmds...)
}

func (m *App) AfterUpdate() tea.Cmd {
	m.list.UpdateProps(m.listProps())
	return nil
}

func (m App) CurrentRoute() reactea.SomeComponent {
	if r, ok := m.routes[reactea.CurrentRoute()]; ok {
		return r
	}
	return nil
}

func (m *App) ToggleItems(items ...int) {
	for _, idx := range items {
		if _, ok := m.Selected[idx]; ok {
			delete(m.Selected, idx)
			m.NumSelected--
		} else if m.NumSelected < m.Limit {
			m.Selected[idx] = struct{}{}
			m.NumSelected++
		}
	}
}

func (c *App) Input(search string) {
	c.inText = search
}

func (c *App) Render(width, height int) string {
	w := c.Width()
	h := c.Height()

	var view []string

	if head := c.renderHeader(w, h); head != "" {
		h -= lipgloss.Height(head)
		view = append(view, head)
	}

	footer := c.renderFooter(w, h)
	if footer != "" {
		h -= lipgloss.Height(footer)
	}

	list := c.list.Render(w, h)
	view = append(view, list)

	if footer != "" {
		view = append(view, footer)
	}

	return lipgloss.JoinVertical(lipgloss.Left, view...)
}

func (c App) renderHeader(w, h int) string {
	var header string
	if c.header != "" {
		header = c.Style.Header.Render(c.header)
	}
	switch reactea.CurrentRoute() {
	case "filter":
		header = c.mainRouter.Render(w, h)
	}
	return header
}

func (c App) renderFooter(w, h int) string {
	var footer string
	if c.footer != "" {
		footer = c.Style.Header.Render(c.footer)
	}
	switch reactea.CurrentRoute() {
	case "status":
		footer = c.mainRouter.Render(w, h)
	}
	return footer
}

func (c *App) ChangeRoute(r string) {
	if p := reactea.CurrentRoute(); p == "" {
		c.PrevRoute = "default"
	} else {
		c.PrevRoute = p
	}
	reactea.SetCurrentRoute(r)
	//c.routes[c.PrevRoute].Destroy()
}

func (m App) Chosen() []map[string]string {
	var chosen []map[string]string
	if len(m.Selected) > 0 {
		for k := range m.Selected {
			chosen = append(chosen, m.Choices[k])
		}
	}
	return chosen
}

func Filter(search string, choices []map[string]string) []item.Item {
	c := item.Choices(choices)
	return c.Filter(search)
}
