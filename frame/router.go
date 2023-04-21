package frame

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/keys"
)

type Router struct {
	*router.Component
	PrevRoute string
}

type Route interface {
	Initialize() router.RouteInitializer
}

func NewRouter() *Router {
	return &Router{
		Component: router.New(),
	}
}

func (c *Router) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	switch msg := msg.(type) {
	case keys.ReturnToListMsg:
		return keys.ChangeRoute("default")

	case keys.ChangeRouteMsg:
		route := msg.Name

		if reactea.CurrentRoute() == route {
			return nil
		}

		switch route {
		case "prev":
			route = c.PrevRoute
		}

		c.PrevRoute = reactea.CurrentRoute()
		reactea.SetCurrentRoute(route)
		return nil
	}

	return c.Component.Update(msg)
}
