package app

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
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
	header      *header.Component

	PrevRoute string
	routes    map[string]reactea.SomeComponent

	Viewport viewport.Model
}

func New(choices []string) *App {
	a := &App{
		width:    util.TermWidth(),
		height:   10,
		Choices:  item.MapChoices(choices),
		Style:    style.DefaultAppStyle(),
		Selected: make(map[int]struct{}),
		Limit:    1,
		routes:   make(map[string]reactea.SomeComponent),
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
	c.input = input.New()
	c.input.Init(input.Props{
		Filter: c.Input,
	})
	c.header = &header.Component{}

	c.routes["list"] = c.list
	c.routes["filter"] = c.input
	c.routes["header"] = c.header
	return nil
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	if reactea.CurrentRoute() == "" {
		reactea.SetCurrentRoute("list")
	}
	reactea.AfterUpdate(c)
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case message.StopFilteringMsg:
		c.Input("")
		c.list.SetKeyMap(keys.VimListKeyMap())
		cmds = append(cmds, message.ChangeRoute("list"))

	case message.StartFilteringMsg:
		c.list.SetKeyMap(keys.DefaultListKeyMap())
		cmds = append(cmds, message.ChangeRoute("filter"))

	case message.ConfirmMsg:
		//c.ConfirmAction = ""
	case message.GetConfirmationMsg:
		//c.ConfirmAction = msg.Question
	case message.ChangeRouteMsg:
		route := msg.Name
		switch route {
		case "header":
			p := header.Props("poot")
			c.header.Init(p)
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
			cmds = append(cmds, message.ChangeRoute("header"))
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
		cmds = append(cmds, c.input.Update(msg))
		c.search = c.inText
	case "help":
	}

	cmds = append(cmds, c.list.Update(msg))

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
	c.inText = search
}

func (c *App) SetContent(lines string) {
	c.Viewport.SetContent(lines)
}

func (c *App) Render(width, height int) string {
	view := c.list.Render(c.Width(), c.Height())

	var header string
	if h := c.header.Render(width, height); h != "" {
		header = c.Style.Header.Render(h)
	}

	switch reactea.CurrentRoute() {
	case "filter":
		header = c.input.Render(c.Width(), c.Height())
	}

	if header != "" {
		view = lipgloss.JoinVertical(lipgloss.Left, header, view)
	}

	//if c.header != "" {
	//view = lipgloss.JoinVertical(lipgloss.Left, c.header, view)
	//}

	//if c.ConfirmAction != "" {
	//c.Footer(fmt.Sprintf("%s\n", c.Style.Confirm.Render(c.ConfirmAction+"(y/n)")))
	//}

	//if c.footer != "" {
	//view = lipgloss.JoinVertical(lipgloss.Left, view, c.footer)
	//}

	return view
	//return lipgloss.JoinVertical(lipgloss.Left, view, reactea.CurrentRoute())
}

func (c *App) ChangeRoute(r string) {
	if p := reactea.CurrentRoute(); p == "" {
		c.PrevRoute = "default"
	} else {
		c.PrevRoute = p
	}
	reactea.SetCurrentRoute(r)
	c.routes[c.PrevRoute].Destroy()
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
