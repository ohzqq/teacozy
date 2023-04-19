package frame

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/util"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	choices item.Choices
}

func NewApp(c []string) *App {
	return &App{
		mainRouter: router.New(),
		choices:    item.SliceToChoices(c),
	}
}

func (m *App) Init(reactea.NoProps) tea.Cmd {
	return m.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := NewComponent()
			return component, component.Init(NewProps(m.choices))
		},
	})
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}
	cmd = c.mainRouter.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *App) Render(w, h int) string {
	return c.mainRouter.Render(util.TermWidth(), 10)
}
