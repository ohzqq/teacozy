package body

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
)

type Component struct {
	reactea.BasicComponent
	SetValue func(string)
}

func New(fn func(string)) *Component {
	c := &Component{
		SetValue: fn,
	}
	return c
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	//var cmd tea.Cmd
	//var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
		//c.Props().SetValue("poot")
		if msg.String() == "b" {
			//c.Props().SetValue("alt body")
			//reactea.SetCurrentRoute("alt body")
		}

	}

	//cmd = c.Props().Component.Update(msg)
	//cmds = append(cmds, cmd)

	//return tea.Batch(cmds...)
	return nil
}

func (c *Component) Render(w, h int) string {
	//view := c.Props().Component.Render(w, h)
	c.SetValue("pooot")

	return "poot"
	//return lipgloss.JoinVertical(lipgloss.Left, view)
}

func (c *Component) Routes() router.Props {
	return map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("body")
		},
		"alt body": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := reactea.Componentify[string](Renderer)
			return comp, comp.Init("alt body")
		},
	}
}

type TestProps = string

func Renderer(p TestProps, w, h int) string {
	return fmt.Sprintf("%s", p)
}
