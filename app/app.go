package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
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
	list        *List
	input       *Input

	Viewport viewport.Model
}

func New(choices []string) *App {
	a := &App{
		mainRouter: router.New(),
		width:      util.TermHeight() - 4,
		height:     util.TermWidth(),
		Choices:    MapChoices(choices),
		Style:      style.DefaultAppStyle(),
		Selected:   make(map[int]struct{}),
		Limit:      10,
	}

	return a
}

func (c *App) listProps() Props {
	p := Props{
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
	c.list = NewList()
	c.list.Init(c.listProps())
	c.input = NewSearch()
	c.input.Init(InputProps{
		Filter: c.Filter,
	})
	return nil
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	if reactea.CurrentRoute() == "" {
		reactea.SetCurrentRoute("list")
	}
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case message.StopFilteringMsg:
		c.Filter("")
		c.list.Init(c.listProps())
		cmds = append(cmds, message.ChangeRoute("list"))

	case message.StartFilteringMsg:
		c.input.Init(InputProps{
			Filter: c.Filter,
		})
		cmds = append(cmds, message.ChangeRoute("filter"))

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
		c.list.Init(c.listProps())
		//case "", "default":
		//NewList().Init(c.listProps())
		//list := c.Viewport.View()
		//default:
		//c.list.Init(c.listProps())
		//cmds = append(cmds, c.mainRouter.Update(msg))
	}

	cmds = append(cmds, c.list.Update(msg))

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
	c.width = width
	c.height = height - 2

	view := c.list.Render(c.Width(), c.Height()-4)

	switch reactea.CurrentRoute() {
	case "filter":
		input := c.input.Render(c.Width(), c.Height())
		view = lipgloss.JoinVertical(lipgloss.Left, input, view)
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

	return lipgloss.JoinVertical(lipgloss.Left, view, reactea.CurrentRoute())
}

type FilterMsg struct {
	Search string
}

func UpdateFilter(search string) tea.Cmd {
	return func() tea.Msg {
		return FilterMsg{Search: search}
	}
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
