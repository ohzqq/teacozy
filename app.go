package teacozy

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/keys"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	router reactea.Component[router.Props]

	Header reactea.SomeComponent
	Body   reactea.SomeComponent
	Footer reactea.SomeComponent

	header string
	body   string
	footer string
}

type Header struct {
	Value string
}

type Footer struct {
	Value string
}

const RoutePlaceholder = ":header/:body/:footer"

func New() *App {
	return &App{
		router: router.New(),
		header: "header",
		body:   "body",
		footer: "footer",
	}
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, c.router.Init(c.Routes()))

	return tea.Batch(cmds...)
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case keys.ChangeRouteMsg:
		//fmt.Println(msg.Name)
		reactea.SetCurrentRoute(msg.Name)
		return nil
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
		if msg.String() == "h" {
			reactea.SetCurrentRoute("header/alt")
		}
		if msg.String() == "b" {
			reactea.SetCurrentRoute("header/body/footer")
		}
		if msg.String() == "f" {
			reactea.SetCurrentRoute("footer/alt")
		}
	}

	//cmds = append(cmds, c.Header.Update(msg))
	//cmds = append(cmds, c.Body.Update(msg))
	//cmds = append(cmds, c.Footer.Update(msg))

	cmds = append(cmds, c.router.Update(msg))

	return tea.Batch(cmds...)
}

func (c *App) Render(w, h int) string {
	header := c.Header.Render(w, h)
	body := c.Body.Render(w, h)
	footer := c.Footer.Render(w, h)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (c *App) Routes() router.Props {
	return map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			c.Header = NewComponent()
			c.Body = NewComponent()
			c.Footer = NewComponent()
			return c, keys.ChangeRoute("header/body/footer")
		},

		RoutePlaceholder: func(params router.Params) (reactea.SomeComponent, tea.Cmd) {
			header := reactea.Componentify[string](Renderer)
			body := reactea.Componentify[string](Renderer)
			footer := reactea.Componentify[string](Renderer)
			var cmds []tea.Cmd
			cmds = append(cmds, header.Init(params["header"]))
			cmds = append(cmds, body.Init(params["body"]))
			cmds = append(cmds, footer.Init(params["footer"]))
			c.Header = header
			c.Body = body
			c.Footer = footer
			return c, tea.Batch(cmds...)
		},
	}
}

type TestProps = string

func Renderer(p TestProps, w, h int) string {
	return fmt.Sprintf("%s", p)
}

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[TestProps]
}

type Props struct {
	Value string
}

func NewComponent() *Component {
	return &Component{}
}

func (c Component) Render(w, h int) string {
	return fmt.Sprintf("%s", c.Props())
}
