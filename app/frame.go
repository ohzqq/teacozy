package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/util"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]
	mainRouter reactea.Component[router.Props]
	*props.Items
	width  int
	height int
	Routes map[string]router.RouteInitializer
}

type Route interface {
	Initializer(*props.Items) router.RouteInitializer
	Name() string
}

func New(choices []map[string]string, routes []Route, opts ...props.Opt) *App {
	app := &App{
		Items:      props.NewItems(choices, opts...),
		mainRouter: router.New(),
		Routes:     make(map[string]router.RouteInitializer),
		width:      util.TermHeight(),
		height:     util.TermWidth(),
	}

	for i, r := range routes {
		name := r.Name()
		if i == 0 {
			app.Routes["default"] = r.Initializer(app.Items)
		}
		app.Routes[name] = r.Initializer(app.Items)
	}

	return app
}

func (c *App) NewProps() *props.Items {
	items := c.Items.Update()
	items.Width = c.width
	items.Height = c.height
	return items
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(c.Routes)
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	c.Snapshot = c.mainRouter.Render(c.Width, c.Height)
	switch msg := msg.(type) {
	case message.ChangeRouteMsg:
		reactea.SetCurrentRoute(msg.Name)
	case message.ReturnSelectionsMsg:
		return reactea.Destroy
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	return c.mainRouter.Update(msg)
}

func (c *App) Render(width, height int) string {
	view := c.mainRouter.Render(width, height)
	return view
}
