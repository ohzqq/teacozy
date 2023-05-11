package teacozy

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/util"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	router *router.Component

	Items       Items
	Selected    map[int]struct{}
	InputValue  string
	currentItem int
	width       int
	height      int

	Pages     map[string]*Page
	Endpoints []string
	page      reactea.SomeComponent
}

func New(routes ...string) *App {
	return &App{
		Pages:     make(map[string]*Page),
		router:    router.New(),
		Endpoints: routes,
		Selected:  make(map[int]struct{}),
		width:     util.TermWidth(),
		height:    util.TermHeight(),
	}
}

func (c *App) AddPage(pages ...*Page) *App {
	for _, p := range pages {
		c.Pages[p.Slug()] = p
	}
	return c
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	routes := map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			return c.Pages["page"], nil
		},
		"": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			return c.Pages["page"], nil
		},
	}

	for _, route := range c.Endpoints {
		routes[filepath.Join(route, ":slug")] = func(params router.Params) (reactea.SomeComponent, tea.Cmd) {
			//fmt.Printf("%+V\n", c.Pages[params["slug"]])
			return c.Pages[params["slug"]], nil
		}
	}

	return c.router.Init(routes)
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
		}
		if msg.String() == "p" {
			reactea.SetCurrentRoute("list/page")
		}
	}

	cmd = c.router.Update(msg)
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
	return c.router.Render(c.width, c.height)
}
