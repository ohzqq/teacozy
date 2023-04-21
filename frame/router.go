// Modified from
// https://github.com/londek/reactea/blob/v0.4.2/router/router.go
// under MIT license
package router

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/keys"
)

type Component struct {
	*router.Component
	PrevRoute string
}

type Props map[string]RouteInitializer
type RouteInitializer func(Params) (reactea.SomeComponent, tea.Cmd)
type Params = map[string]string

func New() *Component {
	return &Component{
		Component: router.New(),
	}
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	switch msg := msg.(type) {
	case keys.ChangeRouteMsg:
		route := msg.Name
		switch route {
		case "prev":
			route = c.PrevRoute
		}
		reactea.SetCurrentRoute(route)
		c.PrevRoute = reactea.CurrentRoute()
		return nil
	}

	return c.Component.Update(msg)
}
