package teacozy

import (
	"os/exec"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/choose"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/field"
	"github.com/ohzqq/teacozy/filter"
	"github.com/ohzqq/teacozy/help"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
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
	body          string
	exec          *exec.Cmd
	execItem      *exec.Cmd
	Style         Style
	help          keys.KeyMap
	Viewport      *viewport.Model
}

type Style struct {
	Confirm lipgloss.Style
	Footer  lipgloss.Style
	Header  lipgloss.Style
}

type Route interface {
	Initializer(*props.Items) router.RouteInitializer
	Name() string
}

func New(props *props.Items, routes []Route) *App {
	app := &App{
		mainRouter: router.New(),
		Routes:     make(map[string]router.RouteInitializer),
		width:      util.TermHeight(),
		height:     util.TermWidth(),
		Style:      DefaultStyle(),
	}
	app.Items = app.NewProps(props)

	if app.Items.Title != "" {
		app.Items.SetHeader(app.Items.Title)
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

func KeymapToProps(km keys.KeyMap) *props.Items {
	p := props.New(props.ChoiceMap(km.Map()), props.Limit(0))
	return p
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

func (c *App) NewProps(props *props.Items) *props.Items {
	c.Footer("")
	props.SetHeader = c.Header
	props.SetBody = c.Body
	props.SetFooter = c.Footer
	props.SetHelp = c.Help
	return props
}

func (c *App) CloneProps() *props.Items {
	items := c.Items.Update()
	items.Width = c.width
	items.Height = c.height
	return items
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	c.Routes["view"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		comp := reactea.Componentify[view.Props](view.Renderer)
		p := view.Props{
			Header: c.header,
			Body:   c.body,
			Footer: c.footer,
		}
		return comp, comp.Init(p)
	}
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
		case "help":
			p := c.NewProps(KeymapToProps(c.help))
			p.Height = c.Items.Height
			p.Width = c.Items.Width
			c.Routes["help"] = help.New().Initializer(p)
		}
		c.PrevRoute = reactea.CurrentRoute()
		reactea.SetCurrentRoute(route)
	case message.ReturnSelectionsMsg:
		switch reactea.CurrentRoute() {
		case "filter":
			if c.HasRoute("choose") {
				reactea.SetCurrentRoute("choose")
			}
		default:
			return reactea.Destroy
		}
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
	c.Body(view)

	if c.header != "" {
		view = lipgloss.JoinVertical(lipgloss.Left, c.header, view)
	}

	//if c.ConfirmAction != "" {
	//c.Footer(fmt.Sprintf("%s\n", c.Style.Confirm.Render(c.ConfirmAction+"(y/n)")))
	//}

	//if c.footer != "" {
	//view = lipgloss.JoinVertical(lipgloss.Left, view, c.footer)
	//}

	vp := viewport.New(width, height)
	vp.SetContent(view)
	//vp.SetContent(
	//  lipgloss.JoinVertical(
	//    lipgloss.Left,
	//    c.header,
	//    c.mainRouter.Render(width, height),
	//    c.footer,
	//  ),
	//)

	return vp.View()
	//view += "\n current " + reactea.CurrentRoute()
	//view += "\n prev " + c.PrevRoute
	//return view
}

func (c *App) Header(f string) {
	c.header = f
}

func (c *App) Body(f string) {
	c.body = f
}

func (c *App) Footer(f string) {
	c.footer = f
}

func (c *App) Help(p keys.KeyMap) {
	c.help = p
}

func Choose(opts ...props.Opt) *App {
	c := choose.New()
	p := props.New(opts...)
	l := New(p, []Route{c})
	pro := reactea.NewProgram(l)
	if err := pro.Start(); err != nil {
		panic(err)
	}
	return l
}

func Form(opts ...props.Opt) *App {
	c := choose.New()
	fi := field.NewField()
	p := props.New(opts...)
	l := New(p, []Route{c, fi})
	pro := reactea.NewProgram(l)
	if err := pro.Start(); err != nil {
		panic(err)
	}
	return l
}

func Filter(opts ...props.Opt) *App {
	c := filter.New()
	p := props.New(opts...)
	l := New(p, []Route{c})
	pro := reactea.NewProgram(l)
	if err := pro.Start(); err != nil {
		panic(err)
	}
	return l
}

func DefaultStyle() Style {
	return Style{
		Confirm: lipgloss.NewStyle().Background(color.Red()).Foreground(color.Black()),
		Footer:  lipgloss.NewStyle().Background(color.Purple()).Foreground(color.Black()),
		Header:  lipgloss.NewStyle().Background(color.Purple()).Foreground(color.Black()),
	}
}
