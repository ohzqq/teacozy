package choose

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/util"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	Choices     []string
	Selected    map[int]struct{}
	numSelected int
	Limit       int
	Width       int
	Height      int
}

func NewRouter(choices ...string) *Component {
	list := &Component{
		Choices:    choices,
		mainRouter: router.New(),
		Height:     4,
		Limit:      2,
		Selected:   make(map[int]struct{}),
	}

	w, h := util.TermSize()
	if list.Height == 0 {
		list.Height = h - 4
	}
	if list.Width == 0 {
		list.Width = w
	}

	return list
}

func (c *Component) NewProps() ChooseProps {
	items := item.New(c.Choices)
	items.Limit = c.Limit
	return ChooseProps{
		Items:      items,
		Selected:   c.Selected,
		ToggleItem: c.ToggleSelection,
	}
}

func (c *Component) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := New()

			return component, component.Init(c.NewProps())
		},
		"filter": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewFilter()

			return component, component.Init(c.NewProps())
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
	return c.mainRouter.Render(c.Width, c.Height)
}

func (m *Component) ToggleSelection(idx int) {
	if _, ok := m.Selected[idx]; ok {
		delete(m.Selected, idx)
		m.numSelected--
	} else if m.numSelected < m.Limit {
		m.Selected[idx] = struct{}{}
		m.numSelected++
	}
}
