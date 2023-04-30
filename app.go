package teacozy

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
)

type App struct {
	reactea.BasicComponent                         // It implements AfterUpdate()
	reactea.BasicPropfulComponent[reactea.NoProps] // It implements props backend - UpdateProps() and Props()

	header reactea.Component[router.Props] // Our router
	body   reactea.Component[router.Props] // Our router
	footer reactea.Component[router.Props] // Our router
}

func New() *App {
	return &App{
		header: router.New(),
		body:   router.New(),
		footer: router.New(),
	}
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, c.header.Init(c.HeaderRoutes()))
	cmds = append(cmds, c.body.Init(c.BodyRoutes()))
	cmds = append(cmds, c.footer.Init(c.FooterRoutes()))

	return tea.Batch(cmds...)
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
		if msg.String() == "a" {
			reactea.SetCurrentRoute("alt")
		}
	}

	cmds = append(cmds, c.header.Update(msg))
	cmds = append(cmds, c.body.Update(msg))
	cmds = append(cmds, c.footer.Update(msg))

	return tea.Batch(cmds...)
}

func (c *App) Render(w, h int) string {
	header := c.header.Render(w, h)
	body := c.body.Render(w, h)
	footer := c.footer.Render(w, h)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (c *App) HeaderRoutes() router.Props {
	return map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("header")
		},
		"alt": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("alt header")
		},
	}
}

func (c *App) BodyRoutes() router.Props {
	return map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("body")
		},
		"alt": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("alt body")
		},
	}
}

func (c *App) FooterRoutes() router.Props {
	return map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("Footer")
		},
		"alt": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("alt footer")
		},
	}
}

type TestProps = string

func Renderer(p TestProps, w, h int) string {
	return fmt.Sprintf("%s", p)
}
