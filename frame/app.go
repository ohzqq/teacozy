package frame

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
	"github.com/ohzqq/teacozy/util"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	mainRouter reactea.Component[router.Props]
	prevRoute  string
	view       reactea.Component[item.Props]

	filter      string
	start       int
	end         int
	width       int
	height      int
	selected    map[int]struct{}
	numSelected int
	limit       int
	cursor      int
	choices     item.Choices
	paginator   *pagy.Paginator
}

func New(c []string) *App {
	a := &App{
		mainRouter: router.New(),
		choices:    item.SliceToChoices(c),
		selected:   make(map[int]struct{}),
		start:      0,
		end:        10,
		cursor:     0,
		view:       reactea.Componentify[item.Props](item.Renderer),
		width:      util.TermWidth(),
		height:     10,
		limit:      10,
	}
	a.paginator = pagy.New(10, len(c))

	return a
}

func (c App) itemProps() item.Props {
	return item.Props{
		Paginator: c.paginator,
		Choices:   c.choices,
		Selected:  c.selected,
		Search:    c.filter,
	}
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.mainRouter.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			component := reactea.Componentify[item.Props](item.Renderer)
			return component, component.Init(c.itemProps())
		},
		"list": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			comp := NewList()
			p := Props{
				Props:       c.itemProps(),
				ToggleItems: c.ToggleItems,
			}
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
	case keys.ChangeRouteMsg:
		route := msg.Name
		c.prevRoute = reactea.CurrentRoute()
		reactea.SetCurrentRoute(route)

	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
		if msg.String() == "n" {
			reactea.SetCurrentRoute("list")
		}
	}

	c.paginator, cmd = c.paginator.Update(msg)
	cmds = append(cmds, cmd)

	cmd = c.mainRouter.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *App) UpdatePagination(p *pagy.Paginator) {
	c.paginator = p
}

func (c *App) Render(w, h int) string {
	view := c.mainRouter.Render(c.width, c.height)
	//view := item.Renderer(c.itemProps(), c.width, c.height)
	//view += fmt.Sprintf("\ncursor %d start %d:end %d", c.cursor, c.start, c.end)
	return view
}

func (c *App) SetCursor(n int) {
	c.cursor = n
}

func (c *App) SetStart(n int) {
	c.start = n
}

func (c *App) SetEnd(n int) {
	c.end = n
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

func (c App) Selected() {
	fmt.Printf("sel %+V\n", c.selected)
}
