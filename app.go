package teacozy

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/choose"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/field"
	"github.com/ohzqq/teacozy/filter"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/util"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	*props.Items

	mainRouter    reactea.Component[router.Props]
	width         int
	height        int
	Routes        map[string]router.RouteInitializer
	ConfirmAction string
	PrevRoute     string
	footer        string
	ConfirmStyle  lipgloss.Style
}

type Route interface {
	Initializer(*props.Items) router.RouteInitializer
	Name() string
}

func New(choices []map[string]string, routes []Route, opts ...props.Opt) *App {
	app := &App{
		Items:        props.New(choices, opts...),
		mainRouter:   router.New(),
		Routes:       make(map[string]router.RouteInitializer),
		width:        util.TermHeight(),
		height:       util.TermWidth(),
		ConfirmStyle: lipgloss.NewStyle().Background(color.Red()).Foreground(color.Black()),
	}
	app.Items.Footer = app.Footer

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
	var cmds []tea.Cmd
	c.Snapshot = c.mainRouter.Render(c.Width, c.Height)
	switch msg := msg.(type) {
	case message.ConfirmMsg:
		c.ConfirmAction = ""
	case message.GetConfirmationMsg:
		c.ConfirmAction = msg.Question
	case message.ChangeRouteMsg:
		route := msg.Name
		switch route {
		case "prev":
			route = c.PrevRoute
		}
		c.PrevRoute = reactea.CurrentRoute()
		reactea.SetCurrentRoute(route)
	case message.ReturnSelectionsMsg:
		return reactea.Destroy
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return reactea.Destroy
		case "y":
			cmds = append(cmds, message.ConfirmCmd(true))
		case "n":
			cmds = append(cmds, message.ConfirmCmd(false))
		}
	}

	cmds = append(cmds, c.mainRouter.Update(msg))
	return tea.Batch(cmds...)
}

func (c *App) Render(width, height int) string {
	view := c.mainRouter.Render(width, height)

	if c.ConfirmAction != "" {
		c.Footer(fmt.Sprintf("%s\n", c.ConfirmStyle.Render(c.ConfirmAction+"(y/n)")))
	}

	if c.footer != "" {
		view = lipgloss.JoinVertical(lipgloss.Left, view, c.footer)
	}

	//view += "\n current " + reactea.CurrentRoute()
	//view += "\n prev " + c.PrevRoute
	return view
}

func (c *App) Footer(f string) {
	c.footer = f
}

func Choose(choices []map[string]string, opts ...props.Opt) *App {
	c := choose.NewChoice()
	l := New(choices, []Route{c}, opts...)
	pro := reactea.NewProgram(l)
	if err := pro.Start(); err != nil {
		panic(err)
	}
	return l
}

func Form(choices []map[string]string, opts ...props.Opt) *App {
	c := choose.NewChoice()
	fi := field.NewField()
	l := New(choices, []Route{c, fi}, opts...)
	pro := reactea.NewProgram(l)
	if err := pro.Start(); err != nil {
		panic(err)
	}
	return l
}

func Filter(choices []map[string]string, opts ...props.Opt) *App {
	c := filter.NewFilter()
	l := New(choices, []Route{c}, opts...)
	pro := reactea.NewProgram(l)
	if err := pro.Start(); err != nil {
		panic(err)
	}
	return l
}
