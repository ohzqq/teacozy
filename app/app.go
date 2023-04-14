package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/app/confirm"
	"github.com/ohzqq/teacozy/app/input"
	"github.com/ohzqq/teacozy/app/item"
	"github.com/ohzqq/teacozy/app/list"
	"github.com/ohzqq/teacozy/app/status"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	Style     style.App
	width     int
	height    int
	footer    string
	header    string
	status    string
	PrevRoute string

	inputValue string
	search     string

	confirm confirm.Props

	list        *list.Component
	Choices     []map[string]string
	Selected    map[int]struct{}
	NumSelected int
	Limit       int
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

func (c *App) Init(reactea.NoProps) tea.Cmd {
	c.list = list.New()
	c.list.SetKeyMap(keys.VimListKeyMap())
	c.list.Init(c.listProps())

	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := new(struct {
				reactea.BasicComponent
				reactea.InvisibleComponent
			})
			return component, nil
		},
		"status": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := status.New()
			return component, component.Init(status.Props{Msg: c.status})
		},
		"filter": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := input.New()
			c.list.SetKeyMap(keys.DefaultListKeyMap())
			return component, component.Init(input.Props{Filter: c.Input})
		},
		"confirm": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := confirm.New()
			p := c.confirm
			if p.Action == nil {
				p.Action = component.Confirmed
			}
			return component, component.Init(p)
		},
	})
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	if reactea.CurrentRoute() == "" {
		reactea.SetCurrentRoute("list")
	}

	reactea.AfterUpdate(c)

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case status.StatusMsg:
		c.SetStatus(msg.Status)
		cmds = append(cmds, message.ChangeRoute("status"))

	case message.StopFilteringMsg:
		c.Input("")
		c.list.SetKeyMap(keys.VimListKeyMap())
		cmds = append(cmds, message.ChangeRoute("list"))
	case message.StartFilteringMsg:
		cmds = append(cmds, message.ChangeRoute("filter"))

	case confirm.GetConfirmationMsg:
		c.confirm = msg.Props
		cmds = append(cmds, message.ChangeRoute("confirm"))
	case confirm.NotConfirmedMsg:
		c.SetFooter("")
		cmds = append(cmds, message.ChangeRoute("list"))

	case message.ChangeRouteMsg:
		route := msg.Name
		switch route {
		case "list":
			c.list.SetCursor(0)
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
			cmds = append(cmds, confirm.Confirm("poot?", status.StatusUpdate("mew")))
		case "/":
			cmds = append(cmds, message.StartFiltering())
		case "ctrl+c":
			return reactea.Destroy
		}
	}
	switch reactea.CurrentRoute() {
	case "filter":
		c.search = c.inputValue
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
	c.inputValue = search
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
	case "confirm":
		fallthrough
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

func (c App) Height() int {
	return c.height
}

func (c App) Width() int {
	return c.width
}

func (c *App) SetHeader(h string) {
	c.header = h
}

func (c *App) SetFooter(h string) {
	c.footer = h
}

func (c *App) SetStatus(h string) {
	c.status = h
}
