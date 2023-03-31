package frame

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/choose"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
)

type Frame struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]
	mainRouter reactea.Component[router.Props]
	*props.Items
	width  int
	height int
	Routes map[string]router.RouteInitializer
	Routez []Route
}

type Route interface {
	Initializer(*props.Items) router.RouteInitializer
	Name() string
}

func New(choices []map[string]string, opts ...props.Opt) *Frame {
	app := &Frame{
		Items:      props.NewItems(choices, opts...),
		mainRouter: router.New(),
		Routes:     make(map[string]router.RouteInitializer),
	}
	app.Routes["default"] = choose.RouteInitializer(choose.Props{Items: app.Items})

	return app
}

func NewFrame(choices []map[string]string, routes []Route, opts ...props.Opt) *Frame {
	app := &Frame{
		Items:      props.NewItems(choices, opts...),
		mainRouter: router.New(),
		Routes:     make(map[string]router.RouteInitializer),
		Routez:     routes,
		width:      50,
		height:     5,
	}
	for _, r := range routes {
		app.Routes[r.Name()] = r.Initializer(app.Items)
	}

	return app
}

func (c *Frame) NewProps() *props.Items {
	items := c.Items.Update()
	items.Width = c.width
	items.Height = c.height
	return items
}

func (c *Frame) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(c.Routes)
}

func (c Frame) routes() map[string]router.RouteInitializer {
	routes := make(map[string]router.RouteInitializer)
	for _, r := range c.Routez {
		routes[r.Name()] = r.Initializer(c.NewProps())
	}
	return routes
}

func (c *Frame) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
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

func (c *Frame) Render(width, height int) string {
	view := c.mainRouter.Render(c.width, c.height)
	return view
}
