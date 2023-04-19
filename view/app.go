package view

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
	*Model
}

func NewApp(c []string) *App {
	return &App{
		mainRouter: router.New(),
		choices:    item.SliceToChoices(c),
		Model:      NewView(),
	}
}

func (m App) itemProps() item.Props {
	p := item.Props{
		Choices:  m.choices,
		Selected: make(map[int]struct{}),
		Start:    0,
		End:      10,
		Cursor:   0,
		//Start:      m.Model.Props().Start,
		//End:        m.Model.Props().End,
		//Cursor:     m.Model.Props().Cursor,
		//Selected:   m.Model.Props().Selected,
		//Selectable: m.Model.Props().Selectable,
	}
	return p
}

func (m *App) Init(reactea.NoProps) tea.Cmd {
	return m.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := item.NewList()
			return component, component.Init(m.itemProps())
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
	c.Model, cmd = c.Model.Update(msg)
	cmds = append(cmds, cmd)

	cmd = c.mainRouter.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *App) Render(w, h int) string {
	return c.mainRouter.Render(util.TermWidth(), 10)
}
