package teacozy

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	Items       Items
	Selected    map[int]struct{}
	InputValue  string
	currentItem int
	width       int
	height      int

	Pages  map[string]*Page
	Routes []string
	page   *Page
}

func New(pages map[string]*Page, routes ...string) *App {
	r := make([]string, len(routes))
	for i, route := range routes {
		r[i] = filepath.Join(route, ":slug")
	}
	return &App{
		Pages:    pages,
		Routes:   r,
		Selected: make(map[int]struct{}),
	}
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.initializePage()
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
		if msg.String() == "s" {
			reactea.SetCurrentRoute("list/slice")
			return nil
		}
		if msg.String() == "p" {
			reactea.SetCurrentRoute("list/page")
			return nil
		}
	}

	cmd = c.page.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c *App) AfterUpdate() tea.Cmd {
	if !reactea.WasRouteChanged() {
		return nil
	}

	if c.page != nil {
		c.page.Destroy()
	}

	c.page = nil

	return c.initializePage()
}

func (c *App) initializePage() tea.Cmd {
	for _, ph := range c.Routes {
		if params, ok := reactea.RouteMatchesPlaceholder(reactea.CurrentRoute(), ph); ok {
			fmt.Println(params)
			if page, ok := c.Pages[params["slug"]]; ok {
				c.page = page
				return nil
			}
		}
	}
	return nil
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
	return c.page.Render(w, h)
}
