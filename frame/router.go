package frame

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
)

type Router struct {
	*router.Component
	PrevRoute    string
	UpdateRoutes func(Route)
}

type Route interface {
	Initialize(*App)
	Initializer(teacozy.Props) router.RouteInitializer
	Name() string
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

	case ChangeRouteMsg:
		c.UpdateRoutes(msg.Route)
		return keys.ChangeRoute(msg.Route.Name())

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
		return keys.UpdateStatus(route)
		//return nil
	}

	return c.Component.Update(msg)
}

type ChangeRouteMsg struct {
	Route Route
}

func ChangeRoute(r Route) tea.Cmd {
	return func() tea.Msg {
		return ChangeRouteMsg{Route: r}
	}
}
