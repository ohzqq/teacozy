// Heavily borrowed from
// https://github.com/londek/reactea/blob/v0.4.2/router/router.go
// under MIT license
package frame

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/keys"
)

type Router struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[RouterProps]

	lastComponent reactea.SomeComponent
	PrevRoute     string
}

type RouterProps map[string]RouteInitializer
type RouteInitializer func(Params) (reactea.SomeComponent, tea.Cmd)
type Params = map[string]string

func NewRouter() *Router {
	return &Router{}
}

func (c *Router) Init(props RouterProps) tea.Cmd {
	c.UpdateProps(props)

	return c.initializeRoute()
}

func (c *Router) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	if c.lastComponent == nil {
		return nil
	}

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

	return c.lastComponent.Update(msg)
}

func (c *Router) AfterUpdate() tea.Cmd {
	// If last route != currentRoute we want to reinitialize the component
	if !reactea.WasRouteChanged() {
		return nil
	}

	if c.lastComponent != nil {
		c.lastComponent.Destroy()
	}

	c.lastComponent = nil

	return c.initializeRoute()
}

func (c *Router) Render(width, height int) string {
	if c.lastComponent != nil {
		return c.lastComponent.Render(width, height)
	}

	return fmt.Sprintf("Couldn't route for \"%s\"", reactea.CurrentRoute())
}

func (c *Router) initializeRoute() tea.Cmd {
	var cmd tea.Cmd

	if initializer, ok := c.Props()[reactea.CurrentRoute()]; ok {
		c.lastComponent, cmd = initializer(nil)
	} else if initializer, params, ok := c.findMatchingRouteInitializer(); ok {
		c.lastComponent, cmd = initializer(params)
	} else if initializer, ok := c.Props()["default"]; ok {
		c.lastComponent, cmd = initializer(nil)
	}

	return cmd
}

func (c *Router) findMatchingRouteInitializer() (RouteInitializer, Params, bool) {
	currentRoute := reactea.CurrentRoute()

	for placeholder, initializer := range c.Props() {
		if params, ok := reactea.RouteMatchesPlaceholder(currentRoute, placeholder); ok {
			return initializer, params, true
		}
	}

	return nil, nil, false
}
