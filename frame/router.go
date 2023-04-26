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

	reactea.BasicPropfulComponent[RouterProps]

	PrevRoute    string
	defaultRoute string
	UpdateRoutes func(Route) tea.Cmd
}

type Route interface {
	Initializer(teacozy.Props) router.RouteInitializer
	Name() string
}

type RouterProps struct {
	Routes      router.Props
	Default     string
	ChangeRoute func(Route) tea.Cmd
}

func NewRouter() *Router {
	return &Router{
		Component: router.New(),
	}
}

func (c *Router) Init(props RouterProps) tea.Cmd {
	c.UpdateProps(props)

	return c.Component.Init(c.Props().Routes)
}

func (c *Router) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	switch msg := msg.(type) {
	case keys.ReturnToListMsg:
		return keys.ChangeRoute("default")

	case ChangeRouteMsg:
		if _, ok := c.Component.Props()[msg.Route.Name()]; ok {
			return keys.ChangeRoute(msg.Route.Name())
		}
		return c.UpdateRoutes(msg.Route)

	case keys.ChangeRouteMsg:
		route := msg.Name

		if reactea.CurrentRoute() == route {
			return nil
		}

		if route == "prev" {
			route = c.PrevRoute
		}

		if route == "default" && c.Props().Default != "default" {
			route = c.Props().Default
		}

		c.PrevRoute = reactea.CurrentRoute()
		reactea.SetCurrentRoute(route)
		//u := fmt.Sprintf("cur %s prev %s", reactea.CurrentRoute(), c.PrevRoute)
		//u := fmt.Sprintf("routes %v", c.Props().Default)
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
