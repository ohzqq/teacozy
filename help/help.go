package help

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
)

type Component struct {
	KeyMap keys.KeyMap
}

func New(km keys.KeyMap) *Component {
	return &Component{KeyMap: km}
}

func (c *Component) Initializer(props teacozy.Props) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := reactea.Componentify[teacozy.Props](c.Render)
		props.Items = c.KeyMap
		props.ReadOnly = true
		props.KeyMap.AddBinds(keys.Esc().AddKeys("q").Cmd(keys.ChangeRoute("prev")))
		return component, component.Init(props)
	}
}

func (c Component) Name() string {
	return "help"
}

func (m *Component) Render(props teacozy.Props, w, h int) string {
	m.view.SetWidth(w)
	m.view.SetHeight(h)
	return m.view.View()
}
