package app

import (
	"path/filepath"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/cmpnt"
	"github.com/ohzqq/teacozy/keys"
	"golang.org/x/exp/slices"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	router    *router.Component
	pages     map[string]*teacozy.Page
	routes    map[string]router.RouteInitializer
	endpoints []string
	prevRoute string
	selected  map[int]struct{}
}

func New(choices teacozy.Items) *App {
	c := &App{
		router:    router.New(),
		pages:     make(map[string]*teacozy.Page),
		routes:    make(map[string]router.RouteInitializer),
		prevRoute: "default",
		selected:  make(map[int]struct{}),
	}
	help := teacozy.NewPage("help")
	help.InitFunc(cmpnt.NewHelp)
	c.pages["help"] = help

	c.NewPage("list", choices)
	//c.pages["list"] = c.pages["default"]

	return c
}

func (c *App) AddPage(page *teacozy.Page) *App {
	c.endpoints = append(c.endpoints, page.Endpoint)
	c.pages[page.Endpoint] = page
	if page.Endpoint != "default" {
		route := filepath.Join(page.Endpoint, ":id")
		c.routes[route] = c.initRoute(page.Endpoint)
	}
	c.pages["help"].AddItems(teacozy.MapToChoices(page.KeyMap().Map()))
	return c
}

func (c *App) NewPage(endpoint string, data ...teacozy.Items) {
	page := cmpnt.NewPage(endpoint, data...)
	c.AddPage(page)
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	c.routes["default"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		return c.pages["list"].Update()
	}
	c.routes["help/:id"] = func(params router.Params) (reactea.SomeComponent, tea.Cmd) {
		idx := slices.Index(c.endpoints, c.pages[params["id"]].Endpoint)
		if idx < 0 {
			idx = 0
		}
		c.pages["help"].SetCurrentPage(idx)
		return c.pages["help"].Update()
	}
	return c.router.Init(c.routes)
}

func (c App) initRoute(endpoint string) router.RouteInitializer {
	return func(params router.Params) (reactea.SomeComponent, tea.Cmd) {
		idx, err := strconv.Atoi(params["id"])
		if err != nil {
			idx = 0
		}
		return c.pages[endpoint].SetCurrentPage(idx).Update()
	}
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	if reactea.CurrentRoute() == "" {
		reactea.SetCurrentRoute("list/0")
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case keys.ChangeRouteMsg:
		route := msg.Name
		switch route {
		case reactea.CurrentRoute():
			return nil
		case "prev":
			route = c.prevRoute
		default:
			c.prevRoute = reactea.CurrentRoute()
		}
		reactea.SetCurrentRoute(route)

	case keys.ShowHelpMsg:
		//page := filepath.Dir(reactea.CurrentRoute())
		page := c.CurrentEndpoint()
		return keys.ChangeRoute(filepath.Join("help", page))

	case keys.UpdateItemMsg:
		//return msg.Cmd(c.Current())

	//case keys.ToggleItemsMsg, keys.ToggleItemMsg:
	//c.ToggleItems(c.Current())
	//cmds = append(cmds, keys.LineDown)

	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
		if msg.String() == "f1" {
			return keys.ShowHelp
		}
		if msg.String() == "o" {
			return keys.ChangeRoute("list/0")
		}

		//for _, k := range c.KeyMap.Keys() {
		//  if key.Matches(msg, k.Binding) {
		//    cmds = append(cmds, k.TeaCmd)
		//  }
		//}

	}

	cmd = c.router.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c App) CurrentEndpoint() string {
	if dir := filepath.Dir(reactea.CurrentRoute()); dir != "." {
		return dir
	}
	if base := filepath.Base(reactea.CurrentRoute()); slices.Contains(c.endpoints, base) {
		return base
	}
	return "default"
}

func (c *App) Render(w, h int) string {
	return c.router.Render(w, h)
}

func (c *App) Selected() map[int]struct{} {
	return c.selected
}

func (c *App) SelectItem(idx int) {
	c.selected[idx] = struct{}{}
}

func (c *App) DeselectItem(idx int) {
	delete(c.selected, idx)
}

var fields = []map[string]string{
	map[string]string{"Artichoke": "Baking "},
	map[string]string{"Bananas": "Flour"},
	map[string]string{"Sprouts": "Barley"},
	map[string]string{"Bean": "four"},
	map[string]string{"Bitter": "Melon"},
	map[string]string{"Cod": "Orange"},
	map[string]string{"Sugar": "Apple"},
	map[string]string{"Cashews": "Cucumber"},
	map[string]string{"Curry": "Currywurst"},
	map[string]string{"Dill": "Dragonfruit"},
	map[string]string{"Eggs": "Furikake"},
	map[string]string{"Garlic": "Gherkinhree"},
	map[string]string{"Ginger": "Grapefruit"},
	map[string]string{"Hazelnuts": "Horseradish"},
	map[string]string{"Jicama": "Kohlrabi"},
	map[string]string{"Leeks": "four"},
	map[string]string{"Milk": "Molasses"},
	map[string]string{"Muesli": "six"},
	map[string]string{"Nopal": "Nectarine"},
	map[string]string{"Nutella": "Milk"},
	map[string]string{"Oatmeal": "Olives"},
	map[string]string{"Papaya": "Gherkin"},
	map[string]string{"Peppers": "Pickle"},
	map[string]string{"Pineapple": "Plantains"},
	map[string]string{"Pocky": "Quince"},
	map[string]string{"Radish": "Ramps"},
	map[string]string{"Tamarind": "Watermelon"},
	map[string]string{"Weißwurst": "Yams"},
	map[string]string{"Yeast": "Yuzu"},
}
