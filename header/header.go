package header

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

type Props struct {
	SetValue  func(string)
	Component reactea.SomeComponent
}

func New() *Component {
	c := &Component{}
	return c
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
		if msg.String() == "h" {
			c.Props().SetValue("alt body")
		}
	}

	cmd = c.Props().Component.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	view := "headerz"

	return lipgloss.JoinVertical(lipgloss.Left, view)
}

func (c *Component) Routes() router.Props {
	return map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("header")
		},
		"alt header": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("alt header")
		},
	}
}

type TestProps = string

func Renderer(p TestProps, w, h int) string {
	return fmt.Sprintf("%s", p)
}
