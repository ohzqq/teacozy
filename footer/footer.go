package footer

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	routes reactea.Component[router.Props]

	model tea.Model
}

func New() *Component {
	c := &Component{
		routes: router.New(),
	}
	c.model = reactea.New(c)
	return c
}

func (c *Component) Init(reactea.NoProps) tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, c.routes.Init(c.Routes()))

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
		if msg.String() == "f" {
			reactea.SetCurrentRoute("alt footer")
		}
	}

	cmds = append(cmds, c.routes.Update(msg))

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	view := c.routes.Render(w, h)

	return lipgloss.JoinVertical(lipgloss.Left, view)
}

func (c *Component) Routes() router.Props {
	return map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("footer")
		},
		"alt footer": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("alt footer")
		},
	}
}

type TestProps = string

func Renderer(p TestProps, w, h int) string {
	return fmt.Sprintf("%s", p)
}
