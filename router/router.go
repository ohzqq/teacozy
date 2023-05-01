package router

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/keys"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	headers map[string]router.RouteInitializer
	mains   map[string]router.RouteInitializer
	footers map[string]router.RouteInitializer

	Header reactea.SomeComponent
	Main   reactea.SomeComponent
	Footer reactea.SomeComponent
}

//type Route struct {
//Path string
//Init router.RouteInitializer
//}

type Route string

const RoutePlaceholder = ":header/:main/:footer"

// func New(def map[string]router.RouteInitializer) *Component {
// func New() *Component {
func New(header, main, footer router.RouteInitializer) *Component {
	c := &Component{
		headers: make(map[string]router.RouteInitializer),
		mains:   make(map[string]router.RouteInitializer),
		footers: make(map[string]router.RouteInitializer),
	}
	//c.headers["default"] = DefaultInitializer()
	//c.mains["default"] = DefaultInitializer()
	//c.footers["default"] = DefaultInitializer()
	c.headers["default"] = header
	c.mains["default"] = main
	c.footers["default"] = footer
	return c
}

func (c *Component) Init(reactea.NoProps) tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, c.initializeRoutes())

	cmds = append(cmds, keys.ChangeRoute("default/default/default"))

	return tea.Batch(cmds...)
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case keys.ChangeRouteMsg:
		reactea.SetCurrentRoute(msg.Name)
		return nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	if c.Header != nil {
		cmds = append(cmds, c.Header.Update(msg))
	}

	if c.Main != nil {
		cmds = append(cmds, c.Main.Update(msg))
	}

	if c.Footer != nil {
		cmds = append(cmds, c.Footer.Update(msg))
	}

	return tea.Batch(cmds...)
}

func (c *Component) AfterUpdate() tea.Cmd {
	if !reactea.WasRouteChanged() {
		return nil
	}

	if c.Header != nil {
		c.Header.Destroy()
	}
	c.Header = nil

	if c.Main != nil {
		c.Main.Destroy()
	}
	c.Main = nil

	if c.Footer != nil {
		c.Footer.Destroy()
	}
	c.Footer = nil

	return c.initializeRoutes()
}

func (c *Component) initializeRoutes() tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if params, ok := reactea.RouteMatchesPlaceholder(reactea.CurrentRoute(), RoutePlaceholder); ok {
		if header, ok := c.headers[params["header"]]; ok {
			c.Header, cmd = header(params)
			cmds = append(cmds, cmd)
		}

		if main, ok := c.mains[params["main"]]; ok {
			c.Main, cmd = main(params)
			cmds = append(cmds, cmd)
		}

		if footer, ok := c.footers[params["footer"]]; ok {
			c.Footer, cmd = footer(params)
			cmds = append(cmds, cmd)
		}
	} else {
		c.Header, cmd = c.headers["default"](nil)
		cmds = append(cmds, cmd)

		c.Main, cmd = c.mains["default"](nil)
		cmds = append(cmds, cmd)

		c.Footer, cmd = c.footers["default"](nil)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	header := c.Header.Render(w, h)
	main := c.Main.Render(w, h)
	footer := c.Footer.Render(w, h)
	return lipgloss.JoinVertical(lipgloss.Left, header, main, footer)
}

func DefaultInitializer() router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		//header := reactea.Componentify[string](Renderer)
		//return header, header.Init("poot")
		ugh := &struct {
			reactea.BasicComponent
			reactea.InvisibleComponent
		}{}
		return ugh, nil
	}
}

type TestProps = string

func Renderer(p TestProps, w, h int) string {
	return fmt.Sprintf("%s", p)
}

func (c *Component) AddHeader(name string, init router.RouteInitializer) *Component {
	c.headers[name] = init
	return c
}

func (c *Component) AddMain(name string, init router.RouteInitializer) *Component {
	c.mains[name] = init
	return c
}

func (c *Component) AddFooter(name string, init router.RouteInitializer) *Component {
	c.footers[name] = init
	return c
}
