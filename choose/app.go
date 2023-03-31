package choose

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
)

type base struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]
	mainRouter reactea.Component[router.Props]

	*props.Items
	Routes router.Props
	Field  *Field
	*Choose
}

func New(choices []map[string]string, opts ...props.Opt) *base {
	app := &base{
		Items:      props.NewItems(choices, opts...),
		mainRouter: router.New(),
		Routes:     make(router.Props),
		Field:      NewField(),
		Choose:     NewChoice(),
	}
	return app
}

func (c *base) Init(reactea.NoProps) tea.Cmd {
	c.Routes["default"] = c.Initializer(c.Items)
	c.Routes["choose"] = c.Initializer(c.Items)
	c.Routes["editField"] = c.Field.Initializer(c.Items)
	return c.mainRouter.Init(c.Routes)
}

func (c *base) Update(msg tea.Msg) tea.Cmd {
	c.Snapshot = c.mainRouter.Render(c.Width, c.Height)
	switch msg := msg.(type) {
	case message.ChangeRouteMsg:
		reactea.SetCurrentRoute(msg.Name)
	case message.StartEditingMsg:
		return message.ChangeRouteCmd("editField")
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	return c.mainRouter.Update(msg)
}

func (c *base) Render(width, height int) string {
	view := c.mainRouter.Render(width, height)
	return view
}
