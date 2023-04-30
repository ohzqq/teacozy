package router

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]
}

type Route struct {
	Path   string
	Header router.RouteInitializer
	Body   router.RouteInitializer
	Footer router.RouteInitializer
}

const RoutePlaceholder = ":header/:body/:footer"

type Props struct {
	Routes map[string]Route
}

func New() *Component {
	return &Component{}
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	var cmds []tea.Cmd

	if c.Props().Header != nil {
		cmds = append(cmds, c.Props().Header.Init(c.HeaderRoutes()))
	}

	cmds = append(cmds, c.body.Init(c.BodyRoutes()))
	cmds = append(cmds, c.footer.Init(c.FooterRoutes()))

	return tea.Batch(cmds...)
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
		if msg.String() == "h" {
			reactea.SetCurrentRoute("header/alt")
		}
		if msg.String() == "b" {
			reactea.SetCurrentRoute("body/alt")
		}
		if msg.String() == "f" {
			reactea.SetCurrentRoute("footer/alt")
		}
	}

	cmds = append(cmds, c.header.Update(msg))
	cmds = append(cmds, c.body.Update(msg))
	cmds = append(cmds, c.footer.Update(msg))

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	header := c.header.Render(w, h)
	body := c.body.Render(w, h)
	footer := c.footer.Render(w, h)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (c *Component) HeaderRoutes() router.Props {
	return map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("header")
		},
		"header/alt": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("alt header")
		},
	}
}

func (c *Component) BodyRoutes() router.Props {
	return map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("body")
		},
		"body/alt": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("alt body")
		},
	}
}

func (c *Component) FooterRoutes() router.Props {
	return map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("Footer")
		},
		"footer/alt": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("alt footer")
		},
	}
}

type TestProps = string

func Renderer(p TestProps, w, h int) string {
	return fmt.Sprintf("%s", p)
}
