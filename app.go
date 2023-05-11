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

	Pages     []*Page
	Endpoints []string
	page      reactea.SomeComponent
}

func New(routes ...string) *App {
	return &App{
		//Pages:     make(map[string]*Page),
		router:    router.New(),
		Endpoints: routes,
		Selected:  make(map[int]struct{}),
		width:     util.TermWidth(),
		height:    util.TermHeight(),
	}
}

func (c *App) AddPage(pages ...*Page) *App {
	c.Pages = pages
	//for _, p := range pages {
	//  c.Pages[p.Slug()] = p
	//}
	return c
}

var blankPage = &struct {
	reactea.BasicComponent
	reactea.InvisibleComponent
}{}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	routes := map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			if len(c.Pages) > 0 {
				return c.Pages[0], nil
			}
			return blankPage, nil
		},
		"": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			if len(c.Pages) > 0 {
				return c.Pages[0], nil
			}
			return blankPage, nil
		},
	}

	for _, route := range c.Endpoints {
		routes[filepath.Join(route, ":slug")] = func(params router.Params) (reactea.SomeComponent, tea.Cmd) {
			for _, p := range c.Pages {
				if p.Slug() == params["$"] {
					return p, nil
				}
			}
			return nil, nil
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

func (c *App) SetSize(w, h int) *App {
	c.width = w
	c.height = h
	return c
}

func (c *App) SetHeight(h int) *App {
	c.height = h
	return c
}

func (c *App) SetWidth(w int) *App {
	c.width = w
	return c
}
