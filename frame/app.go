package frame

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/util"
	"github.com/ohzqq/teacozy/view"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]

	filter      string
	start       int
	end         int
	selected    map[int]struct{}
	numSelected int
	limit       int
	cursor      int
	choices     item.Choices
}

func New(c []string) *App {
	return &App{
		mainRouter: router.New(),
		choices:    item.SliceToChoices(c),
		selected:   make(map[int]struct{}),
		start:      0,
		end:        10,
		cursor:     0,
	}
}

func (c App) itemProps() item.Props {
	return item.Props{
		Choices:  c.choices,
		Selected: make(map[int]struct{}),
		Start:    0,
		End:      10,
		Cursor:   0,
		Search:   c.filter,
	}
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := reactea.Componentify[item.Props](item.Renderer)
			return component, component.Init(c.itemProps())
		},
		"nav": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := view.New()
			p := view.Props{Props: c.itemProps()}
			return comp, comp.Init(p)
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

func (c *App) SetCursor(n int) {
	c.cursor = n
}

func (m *App) ToggleItems(items ...int) {
	for _, idx := range items {
		if _, ok := m.selected[idx]; ok {
			delete(m.selected, idx)
			m.numSelected--
		} else if m.numSelected < m.limit {
			m.selected[idx] = struct{}{}
			m.numSelected++
		}
	}
}

func (c *App) Filter(search string) []item.Item {
	return c.choices.Filter(search)
}
