package choose

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/props"
)

type base struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]
	mainRouter reactea.Component[router.Props]
	*props.Items
	Routes router.Props

	*Choose
}

func New(choices []map[string]string, opts ...props.Opt) *base {
	app := &base{
		Items:      props.NewItems(choices, opts...),
		mainRouter: router.New(),
		Routes:     make(router.Props),
		Choose:     NewChoice(),
	}
	return app
}

func (c *base) Init(reactea.NoProps) tea.Cmd {
	c.Routes["default"] = c.Initializer(c.Items)
	//c.Routes["editField"] = c.Field.Initializer(c.Items.Current)
	c.Routes["editField"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := NewField()
		return component, component.Init(FieldProps{Item: c.Items.Current})
	}

	return c.mainRouter.Init(c.Routes)
}

func (c *base) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	return c.mainRouter.Update(msg)
}

func (c *base) Render(width, height int) string {
	return c.mainRouter.Render(width, height)
}
