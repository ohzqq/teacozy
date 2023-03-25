package react

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
)

type Component struct {
	reactea.BasicComponent                         // It implements AfterUpdate()
	reactea.BasicPropfulComponent[reactea.NoProps] // It implements props backend - UpdateProps() and Props()

	mainRouter reactea.Component[router.Props] // Our router

	text string // The name
}

func New() *Component {
	return &Component{
		mainRouter: router.New(),
	}
}

func (c *Component) Init(reactea.NoProps) tea.Cmd {
	// Does it remind you of something? react-router!
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewInput()

			return component, component.Init(InputProps{
				SetText: c.setText, // Can also use "lambdas" (function can be created here)
			})
		},
		"displayname": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			// RouteInitializer wants SomeComponent so we have to convert
			// Stateless component (renderer) to Component
			component := reactea.Componentify[string](DisplayRenderer)

			return component, component.Init(c.text)
		},
	})
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	return c.mainRouter.Update(msg)
}

func (c *Component) Render(width, height int) string {
	return c.mainRouter.Render(width, height)
}

func (c *Component) setText(text string) {
	c.text = text
}
