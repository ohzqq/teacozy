package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
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
	PrevRoute   string
	Choices     []map[string]string
	Selected    map[int]struct{}
	NumSelected int
	Limit       int
	search      string

	Viewport viewport.Model
}

func New(choices []string) *App {
	a := &App{
		mainRouter: router.New(),
		width:      util.TermHeight(),
		height:     util.TermWidth(),
		Choices:    MapChoices(choices),
		Style:      style.DefaultAppStyle(),
		Selected:   make(map[int]struct{}),
		Limit:      10,
	}

	a.Viewport = viewport.New(a.width, a.height)

	return a
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewList()

			return component, component.Init(Props{
				Matches:     Filter(c.search, c.Choices),
				Selected:    c.Selected,
				Width:       c.Viewport.Width,
				Height:      c.Viewport.Height,
				ToggleItems: c.ToggleItems,
				SetContent:  c.SetContent,
			})
		},
		"filter": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewSearch()

			return component, component.Init(InputProps{
				Filter: c.Filter,
			})
		},
	})
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case keys.ToggleMsg:

	case message.ConfirmMsg:
		//c.ConfirmAction = ""
	case message.GetConfirmationMsg:
		//c.ConfirmAction = msg.Question
	case message.ChangeRouteMsg:
		route := msg.Name
		switch route {
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
			//if c.HasRoute("choose") {
			//c.ChangeRoute("choose")
			//}
		default:
			return reactea.Destroy
		}
	case message.QuitMsg:
		return reactea.Destroy
	case tea.KeyMsg:
		switch msg.String() {
		case "/":
			reactea.SetCurrentRoute("filter")
		case "ctrl+c":
			return reactea.Destroy
		case "y":
			cmds = append(cmds, message.Confirm(true))
		case "n":
			cmds = append(cmds, message.Confirm(false))
		}
	}

	cmds = append(cmds, c.mainRouter.Update(msg))
	return tea.Batch(cmds...)
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

func (c *App) Filter(search string) {
	c.search = search
}

func (c *App) SetContent(lines string) {
	c.Viewport.SetContent(lines)
}

func (c *App) Render(width, height int) string {
	widget := c.mainRouter.Render(width, height)

	c.Viewport.Width = width
	c.Viewport.Height = height - 2

	view := c.Viewport.View()

	switch reactea.CurrentRoute() {
	case "filter":
		view = lipgloss.JoinVertical(lipgloss.Left, widget, view)
		//view = widget
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
}

func (c *App) ChangeRoute(r string) {
	//if p := reactea.CurrentRoute(); p == "" {
	//c.PrevRoute = "default"
	//} else {
	//c.PrevRoute = p
	//}
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

func Filter(search string, choices []map[string]string) []Item {
	matches := []Item{}
	for i, choice := range choices {
		for label, str := range choice {
			match := NewItem(str, i)
			match.Label = label

			search = strings.ToLower(search)
			matchedString := strings.ToLower(str)

			index := strings.Index(matchedString, search)
			if index >= 0 {
				matchedIndexes := []int{}
				for s := range search {
					matchedIndexes = append(matchedIndexes, index+s)
				}
				match.MatchedIndexes = matchedIndexes
				matches = append(matches, match)
			}
		}
	}
	return matches
}
