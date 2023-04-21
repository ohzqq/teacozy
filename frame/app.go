package frame

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
	"github.com/ohzqq/teacozy/util"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter *Router
	Routes     map[string]router.RouteInitializer

	filter      string
	start       int
	end         int
	width       int
	height      int
	selected    map[int]struct{}
	numSelected int
	limit       int
	cursor      int
	choices     item.Choices
	paginator   *pagy.Paginator
}

type Opt func(*App)

func New(c []string, opts ...Opt) *App {
	a := &App{
		mainRouter: NewRouter(),
		Routes:     make(map[string]router.RouteInitializer),
		choices:    item.SliceToChoices(c),
		selected:   make(map[int]struct{}),
		start:      0,
		end:        10,
		cursor:     0,
		width:      util.TermWidth(),
		height:     10,
		limit:      10,
	}
	a.paginator = pagy.New(10, len(c))

	a.NewRoute(a)
	a.NewRoute(NewList())

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (c App) ItemProps() item.Props {
	return item.Props{
		Paginator: c.paginator,
		Choices:   c.choices,
		Selected:  c.selected,
		Search:    c.filter,
	}
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(c.Routes)
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	switch reactea.CurrentRoute() {
	case "":
		reactea.SetCurrentRoute("default")
	}

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
		if msg.String() == "n" {
			cmds = append(cmds, keys.ChangeRoute("list"))
		}
	}

	c.paginator, cmd = c.paginator.Update(msg)
	cmds = append(cmds, cmd)

	cmd = c.mainRouter.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *App) Render(w, h int) string {
	view := c.mainRouter.Render(c.width, c.height)
	//view := item.Renderer(c.itemProps(), c.width, c.height)
	view += fmt.Sprintf("\ncurrent %v\nprev %v", reactea.CurrentRoute(), c.mainRouter.PrevRoute)
	return view
}

func (c *App) NewRoute(r Route) {
	r.Initialize(c)
}

func (m *App) ToggleItems(items ...int) {
	for _, idx := range items {
		if _, ok := m.selected[idx]; ok {
			delete(m.selected, idx)
			m.numSelected--
		} else if m.numSelected < m.limit {
			m.selected[idx] = struct{}{}
			m.numSelected++
		}
	}
}

func (c *App) Filter(search string) []item.Item {
	return c.choices.Filter(search)
}

func (c App) Selected() {
	fmt.Printf("sel %+V\n", c.selected)
}

func (c *App) Initialize(a *App) {
	a.Routes["default"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := reactea.Componentify[item.Props](item.Renderer)
		return component, component.Init(a.ItemProps())
	}
}
