package router

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/keys"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	headers map[string]RouteInitializer
	mains   map[string]RouteInitializer
	footers map[string]RouteInitializer

	Header reactea.SomeComponent
	Main   reactea.SomeComponent
	Footer reactea.SomeComponent
}

type RouteInitializer func() (reactea.SomeComponent, tea.Cmd)

const RoutePlaceholder = ":header/:main/:footer"

func New(main RouteInitializer) *Component {
	c := &Component{
		headers: make(map[string]RouteInitializer),
		mains:   make(map[string]RouteInitializer),
		footers: make(map[string]RouteInitializer),
	}
	c.headers["default"] = DefaultInitializer()
	c.mains["default"] = main
	c.footers["default"] = DefaultInitializer()

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

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case keys.ChangeRouteMsg:
		reactea.SetCurrentRoute(msg.Name)
		return nil
	case ChangeHeaderMsg:
		if header, ok := c.headers[msg.Name]; ok {
			c.Header, cmd = header()
		} else {
			c.Header, cmd = c.headers["default"]()
		}
		cmds = append(cmds, cmd)
	case ChangeMainMsg:
		if main, ok := c.mains[msg.Name]; ok {
			c.Main, cmd = main()
		} else {
			c.Main, cmd = c.mains["default"]()
		}
		cmds = append(cmds, cmd)
	case ChangeFooterMsg:
		if footer, ok := c.footers[msg.Name]; ok {
			c.Footer, cmd = footer()
		} else {
			c.Footer, cmd = c.footers["default"]()
		}
		cmds = append(cmds, cmd)
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
			c.Header, cmd = header()
		}
		cmds = append(cmds, cmd)

		if main, ok := c.mains[params["main"]]; ok {
			c.Main, cmd = main()
		}
		cmds = append(cmds, cmd)

		if footer, ok := c.footers[params["footer"]]; ok {
			c.Footer, cmd = footer()
		}
		cmds = append(cmds, cmd)
	} else {
		//c.Header, cmd = c.headers["default"]()
		//cmds = append(cmds, cmd)

		//c.Main, cmd = c.mains["default"]()
		//cmds = append(cmds, cmd)

		//c.Footer, cmd = c.footers["default"]()
		//cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	var views []string

	if c.Header != nil {
		if header := c.Header.Render(w, h); header != "" {
			views = append(views, header)
		}
	}

	if c.Main != nil {
		if main := c.Main.Render(w, h); main != "" {
			views = append(views, main)
		}
	}

	if c.Footer != nil {
		if footer := c.Footer.Render(w, h); footer != "" {
			views = append(views, footer)
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func DefaultInitializer() RouteInitializer {
	return func() (reactea.SomeComponent, tea.Cmd) {
		return invisibleComponent(), nil
	}
}

func invisibleComponent() reactea.SomeComponent {
	return &struct {
		reactea.BasicComponent
		reactea.InvisibleComponent
	}{}
}

type TestProps = string

func Renderer(p TestProps, w, h int) string {
	return fmt.Sprintf("%s", p)
}

func (c *Component) AddHeader(name string, init RouteInitializer) *Component {
	c.headers[name] = init
	return c
}

func (c *Component) AddMain(name string, init RouteInitializer) *Component {
	c.mains[name] = init
	return c
}

func (c *Component) AddFooter(name string, init RouteInitializer) *Component {
	c.footers[name] = init
	return c
}

type ChangeHeaderMsg struct {
	Name string
}

func ChangeHeader(name string) tea.Cmd {
	return func() tea.Msg {
		return ChangeHeaderMsg{Name: name}
	}
}

type ChangeMainMsg struct {
	Name string
}

func ChangeMain(name string) tea.Cmd {
	return func() tea.Msg {
		return ChangeMainMsg{Name: name}
	}
}

type ChangeFooterMsg struct {
	Name string
}

func ChangeFooter(name string) tea.Cmd {
	return func() tea.Msg {
		return ChangeFooterMsg{Name: name}
	}
}
