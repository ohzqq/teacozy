package router

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/app/list"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	mainRouter *router.Component
	List       *list.Component
}

type Props struct {
	List   list.Props
	Routes router.Props
}

func New() *Component {
	return &Component{
		List: list.New(),
	}
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	c.List.Init(props.List)
	return c.mainRouter.Init(props.Routes)
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	var cmds []tea.Cmd

	cmds = append(cmds, c.List.Update(msg))
	cmds = append(cmds, c.mainRouter.Update(msg))

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	return ""
}
