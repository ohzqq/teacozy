package teacozy

import (
	"fmt"
	"os/exec"

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
	Style         Style
	help          keys.KeyMap
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
		Items:      props,
		mainRouter: router.New(),
		Routes:     make(map[string]router.RouteInitializer),
		width:      util.TermHeight(),
		height:     util.TermWidth(),
		Style:      DefaultStyle(),
	}
	app.Items.Footer = app.Footer
	app.Items.Header = app.Header
	app.Items.Help = app.Help

	if app.Items.Title != "" {
		app.Items.Header(app.Items.Title)
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
	var c []map[string]string
	for _, k := range km {
		c = append(c, map[string]string{k.Help().Key: k.Help().Desc})
	}
	p := props.New(c, props.Limit(0))
	return p
}

func (c *App) CloneProps() *props.Items {
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
	//case message.ShowHelpMsg:
	//fmt.Println("ehlp")
	//cmds = append(cmds, message.ChangeRouteCmd("help"))
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
			p := KeymapToProps(c.help)
			p.Height = c.Items.Height
			p.Width = c.Items.Width
			c.Footer("")
			c.Routes["help"] = help.New().Initializer(p)
		}
		c.PrevRoute = reactea.CurrentRoute()
		reactea.SetCurrentRoute(route)
	case message.ReturnSelectionsMsg:
		return reactea.Destroy
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return reactea.Destroy
		case "ctrl+h":
			cmds = append(cmds, message.ShowHelpCmd())
			//println("help")
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

	if c.header != "" {
		//h := c.Style.Header.Render(c.header)
		view = lipgloss.JoinVertical(lipgloss.Left, c.header, view)
	}

	if c.ConfirmAction != "" {
		c.Footer(fmt.Sprintf("%s\n", c.Style.Confirm.Render(c.ConfirmAction+"(y/n)")))
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

func (c *App) Header(f string) {
	c.header = f
}

func (c *App) Help(p keys.KeyMap) {
	c.help = p
}

func Choose(choices []map[string]string, opts ...props.Opt) *App {
	c := choose.NewChoice()
	p := props.New(choices, opts...)
	l := New(p, []Route{c})
	pro := reactea.NewProgram(l)
	if err := pro.Start(); err != nil {
		panic(err)
	}
	return l
}

func Form(choices []map[string]string, opts ...props.Opt) *App {
	c := choose.NewChoice()
	fi := field.NewField()
	p := props.New(choices, opts...)
	l := New(p, []Route{c, fi})
	pro := reactea.NewProgram(l)
	if err := pro.Start(); err != nil {
		panic(err)
	}
	return l
}

func Filter(choices []map[string]string, opts ...props.Opt) *App {
	c := filter.NewFilter()
	p := props.New(choices, opts...)
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
