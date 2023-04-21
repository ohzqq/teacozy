package frame

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy"
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
	title       string
	choices     item.Choices
	paginator   *pagy.Paginator
}

type Opt func(*App)

func New(opts ...Opt) *App {
	a := &App{
		mainRouter: NewRouter(),
		Routes:     make(map[string]router.RouteInitializer),
		selected:   make(map[int]struct{}),
		start:      0,
		cursor:     0,
		width:      util.TermWidth(),
		height:     10,
		limit:      10,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (c *App) ItemProps() teacozy.Props {
	if reactea.CurrentRoute() == "default" {
		c.paginator.SetPerPage(c.height)
	}
	props := teacozy.NewProps()
	props.Paginator = c.paginator
	//props.Choices = c.choices
	props.Items = c.choices
	props.Selected = c.selected
	props.Search = c.filter
	return props
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	c.NewRoute(c)
	if c.limit == -1 {
		c.limit = c.choices.Len()
	}
	c.paginator = pagy.New(c.height, c.choices.Len())
	c.paginator.SetKeyMap(DefaultKeyMap())
	c.end = c.height
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
		if msg.String() == "/" {
			cmds = append(cmds, keys.ChangeRoute("filter"))
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
	//view += fmt.Sprintf("\ncurrent %v\nprev %v", reactea.CurrentRoute(), c.mainRouter.PrevRoute)
	return view
}

func (c *App) NewRoute(r Route) {
	r.Initialize(c)
}

func (c *App) SetKeyMap(km keys.KeyMap) *App {
	c.paginator.SetKeyMap(km)
	return c
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

func (c App) Selected() {
	fmt.Printf("sel %+V\n", c.selected)
}

func (c *App) Initialize(a *App) {
	a.Routes["default"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := reactea.Componentify[teacozy.Props](teacozy.Renderer)
		return component, component.Init(a.ItemProps())
	}
}

func WithSlice[E any](c []E) Opt {
	return func(a *App) {
		a.choices = item.SliceToChoices(c)
	}
}

func WithMap[K comparable, V any, M ~map[K]V](c []M) Opt {
	return func(a *App) {
		a.choices = item.MapToChoices(c)
	}
}

func WithRoute(r Route) Opt {
	return r.Initialize
}

func NoLimit() Opt {
	return func(a *App) {
		a.limit = -1
	}
}

func WithLimit(l int) Opt {
	return func(a *App) {
		a.limit = l
	}
}

func WithTitle(t string) Opt {
	return func(a *App) {
		a.title = t
	}
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
