package teacozy

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/choose"
	"github.com/ohzqq/teacozy/field"
	"github.com/ohzqq/teacozy/filter"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
	"github.com/ohzqq/teacozy/view"
	"golang.org/x/exp/slices"
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
	header        string
	exec          *exec.Cmd
	execItem      *exec.Cmd
	Style         style.App
	help          keys.KeyMap
}

type Route func() RouteInitializer

type RouteInitializer interface {
	Initializer(*props.Items) router.RouteInitializer
	Name() string
}

func New(props *props.Items, routes ...Route) *App {
	app := &App{
		mainRouter: router.New(),
		Routes:     make(map[string]router.RouteInitializer),
		width:      util.TermHeight(),
		height:     util.TermWidth(),
		Style:      style.DefaultAppStyle(),
	}
	app.Items = app.NewProps(props)

	if app.Items.Title != "" {
		app.Items.SetHeader(app.Items.Title)
	}

	for i, init := range routes {
		r := init()
		name := r.Name()
		if i == 0 {
			app.Routes["default"] = r.Initializer(app.Items)
		}
		app.Routes[name] = r.Initializer(app.Items)
	}

	return app
}

func (c *App) NewProps(props *props.Items) *props.Items {
	c.Footer("")
	props.SetHeader = c.Header
	props.SetFooter = c.Footer
	props.SetHelp = c.Help
	props.TotalLines = c.TotalLines
	return props
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
		case "filter":
			c.Footer("")
		case "prev":
			route = c.PrevRoute
		case "help":
			p := c.NewProps(KeymapToProps(c.help))
			p.Height = c.Items.Height
			p.Width = c.Items.Width
			c.Routes["help"] = view.New().Initializer(p)
		}
		c.ChangeRoute(route)
	case message.ReturnSelectionsMsg:
		switch reactea.CurrentRoute() {
		case "filter":
			if c.HasRoute("choose") {
				c.ChangeRoute("choose")
			}
		default:
			return reactea.Destroy
		}
	case message.QuitMsg:
		return reactea.Destroy
	case tea.KeyMsg:
		switch msg.String() {
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

func (c *App) Render(width, height int) string {
	view := c.mainRouter.Render(width, height)

	if c.header != "" {
		view = lipgloss.JoinVertical(lipgloss.Left, c.header, view)
	}

	if c.ConfirmAction != "" {
		c.Footer(fmt.Sprintf("%s\n", c.Style.Confirm.Render(c.ConfirmAction+"(y/n)")))
	}

	if c.footer != "" {
		view = lipgloss.JoinVertical(lipgloss.Left, view, c.footer)
	}

	return view
}

func (c *App) TotalLines(f int) {
	c.Items.Lines = f
}

func (c *App) Header(f string) {
	c.header = f
}

func (c *App) Footer(f string) {
	c.footer = f
}

func (c *App) Help(p keys.KeyMap) {
	c.help = p
}

func (c *App) ChangeRoute(r string) {
	if p := reactea.CurrentRoute(); p == "" {
		c.PrevRoute = "default"
	} else {
		c.PrevRoute = p
	}
	reactea.SetCurrentRoute(r)
}

func (c App) ListRoutes() []string {
	var r []string
	for n, _ := range c.Routes {
		r = append(r, n)
	}
	return r
}

func (c App) HasRoute(r string) bool {
	return slices.Contains(c.ListRoutes(), r)
}

func WithChoice() Route {
	return func() RouteInitializer {
		return choose.New()
	}
}

func WithFilter() Route {
	return func() RouteInitializer {
		return filter.New()
	}
}

func WithView() Route {
	return func() RouteInitializer {
		return view.New()
	}
}

func WithForm() Route {
	return func() RouteInitializer {
		return field.New()
	}
}

func Choose(opts ...props.Opt) *App {
	p, err := props.New(opts...)
	if err != nil {
		panic(err)
	}
	l := New(p, WithChoice())
	pro := reactea.NewProgram(l)
	if err := pro.Start(); err != nil {
		panic(err)
	}
	return l
}

func Form(opts ...props.Opt) *App {
	p, err := props.New(opts...)
	if err != nil {
		panic(err)
	}
	l := New(p, WithChoice(), WithForm())
	pro := reactea.NewProgram(l)
	if err := pro.Start(); err != nil {
		panic(err)
	}
	return l
}

func Filter(opts ...props.Opt) *App {
	p, err := props.New(opts...)
	if err != nil {
		panic(err)
	}
	l := New(p, WithFilter())
	pro := reactea.NewProgram(l)
	if err := pro.Start(); err != nil {
		panic(err)
	}
	return l
}

func KeymapToProps(km keys.KeyMap) *props.Items {
	p, _ := props.New(props.ChoiceMap(km.Map()), props.Limit(0))
	return p
}
