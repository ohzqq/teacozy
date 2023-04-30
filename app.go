package teacozy

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/body"
	"github.com/ohzqq/teacozy/footer"
	"github.com/ohzqq/teacozy/header"
)

type App struct {
	reactea.BasicComponent                         // It implements AfterUpdate()
	reactea.BasicPropfulComponent[reactea.NoProps] // It implements props backend - UpdateProps() and Props()

	body   *body.Component
	header *header.Component
	footer *footer.Component
}

func New() *App {
	return &App{
		header: header.New(),
		body:   body.New(),
		footer: footer.New(),
	}
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, c.header.Init(reactea.NoProps{}))
	cmds = append(cmds, c.body.Init(reactea.NoProps{}))
	cmds = append(cmds, c.footer.Init(reactea.NoProps{}))

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

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer, reactea.CurrentRoute())
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
