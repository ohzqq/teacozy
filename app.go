package teacozy

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	Items       Items
	Selected    map[int]struct{}
	Search      string
	currentItem int

	Routes map[PlaceHolder]PageComponent
	page   *Page
}

func New(routes map[PlaceHolder]PageComponent) *App {
	return &App{
		Routes:   routes,
		Selected: make(map[int]struct{}),
	}
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.InitializePage()
}

func (c *App) InitializePage() tea.Cmd {
	for ph, page := range c.Routes {
		if params, ok := ph.Matches(reactea.CurrentRoute()); ok {
			p, cmd := page.Initialize(params)
			c.page = p
			return cmd
		}
	}
	return nil
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
	}

	c.page, cmd = c.page.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *App) SetCurrentItem(idx int) {
	c.currentItem = idx
}

func (c *App) CurrentItem() int {
	return c.currentItem
}

func (c *App) Render(w, h int) string {
	if c.page == nil {
		return "404 not found"
	}
	return c.page.View()
}
