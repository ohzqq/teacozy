package list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
)

type List struct {
	reactea.BasicComponent                         // It implements AfterUpdate()
	reactea.BasicPropfulComponent[reactea.NoProps] // It implements props backend - UpdateProps() and Props()

	mainRouter reactea.Component[router.Props] // Our router

	Choices     []string
	choiceMap   []map[string]string
	Items       []Item
	Matches     []Item
	Selected    map[int]struct{}
	Limit       int
	numSelected int
	Cursor      int
}

func New(items ...string) *List {
	return &List{
		mainRouter: router.New(),
	}
}

//func (c *List) Init(reactea.NoProps) tea.Cmd {
//  // Does it remind you of something? react-router!
//  return c.mainRouter.Init(map[string]router.RouteInitializer{
//    "default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
//      component := input.New()

//      return component, component.Init(input.Props{
//        SetText: c.setText, // Can also use "lambdas" (function can be created here)
//      })
//    },
//    "displayname": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
//      // RouteInitializer wants SomeComponent so we have to convert
//      // Stateless component (renderer) to Component
//      component := reactea.Componentify[string](displayname.Renderer)

//      return component, component.Init(c.text)
//    },
//  })
//}

func (c *List) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	return c.mainRouter.Update(msg)
}

func (c *List) Render(width, height int) string {
	return c.mainRouter.Render(width, height)
}
