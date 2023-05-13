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

	router      *router.Component
	pages       map[string]*Page
	endpoints   []string
	prevRoute   string
	CurrentItem int
	Selected    map[int]struct{}

	teacozy.State
}

func New(choices teacozy.Items) *App {
	c := &App{
		router:    router.New(),
		pages:     make(map[string]*Page),
		prevRoute: "default",
	}
	c.NewPage("default", choices)
	c.NewPage("help", teacozy.MapToChoices(c.pages["default"].KeyMap().Map()))
	c.pages["help"].Initializer = cmpnt.NewHelp

	return c
}

func (c *App) NewPage(endpoint string, data ...teacozy.Items) {
	c.endpoints = append(c.endpoints, endpoint)
	props := cmpnt.PagerProps{
		SetCurrent: c.SetCurrent,
		Current:    c.Current,
	}
	if len(data) > 0 {
		props.Items = data[0]
	}
	c.pages[endpoint] = NewPage(endpoint, data...)
	c.pages[endpoint].Init(props)
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.router.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			return c.pages["default"].UpdateProps(""), nil
		},
		"help/:name": func(params router.Params) (reactea.SomeComponent, tea.Cmd) {
			idx := slices.Index(c.endpoints, c.pages[params["name"]].Name)
			page := c.pages["help"].UpdateProps(strconv.Itoa(idx))
			return page, nil
		},
	})
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)

	if reactea.CurrentRoute() == "" {
		reactea.SetCurrentRoute("default")
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
	//cmds = append(cmds, keys.ChangeRoute("help"))

	case keys.ShowHelpMsg:
		page := filepath.Base(reactea.CurrentRoute())
		return keys.ChangeRoute(filepath.Join("help", page))

	case keys.UpdateItemMsg:
		return msg.Cmd(c.Current())

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

func (c *App) Render(w, h int) string {
	return c.router.Render(w, h)
}

func (c App) Current() int {
	return c.CurrentItem
}

func (m *App) SetCurrent(idx int) {
	m.CurrentItem = idx
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
