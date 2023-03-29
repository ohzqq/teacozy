package list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type List struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	Choices     []string
	Items       []Item
	choiceMap   []map[string]string
	Selected    map[int]struct{}
	Matches     []Item
	Limit       int
	numSelected int
	Cursor      int
	quitting    bool
	header      string
	Placeholder string
	Prompt      string
	Width       int
	Height      int
	Style       style.List
}

func New(items ...string) *List {
	list := &List{
		Items:      ChoicesToMatch(items),
		Choices:    items,
		Selected:   make(map[int]struct{}),
		mainRouter: router.New(),
		Limit:      1,
		Height:     10,
	}
	list.Matches = list.Items

	w, h := util.TermSize()
	if list.Height == 0 {
		list.Height = h - 4
	}
	if list.Width == 0 {
		list.Width = w
	}

	return list
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
