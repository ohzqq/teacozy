package app

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/cmpnt"
	"github.com/ohzqq/teacozy/keys"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	router    *router.Component
	pages     map[string]*cmpnt.Pager
	endpoints []string

	*cmpnt.Pager
}

func New(opts ...cmpnt.Option) *App {
	c := &App{
		router: router.New(),
		endpoints: []string{
			"main/:name",
			"help/:name",
		},
		pages: make(map[string]*cmpnt.Pager),
	}
	c.Pager = cmpnt.New(opts...)
	c.pages["default"] = c.Pager

	return c
}

func (c *App) Init(reactea.NoProps) tea.Cmd {
	return c.router.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			return c.pages["default"], nil
		},
		"help/:name": func(params router.Params) (reactea.SomeComponent, tea.Cmd) {
			if p, ok := c.pages[params["name"]]; ok {
				opts := []cmpnt.Option{
					cmpnt.WithMap(p.KeyMap().Map()),
					cmpnt.ReadOnly(),
				}
				page := cmpnt.New(opts...)
				return page, nil
			}
			return c.pages["default"], nil
		},
	})
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
	if reactea.CurrentRoute() == "" {
		reactea.SetCurrentRoute("default")
	}
	reactea.AfterUpdate(c)
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case keys.ShowHelpMsg:
		cmds = append(cmds, keys.ChangeRoute("help"))
		//help := NewProps(c.help)
		//help.SetName("help")
		//return ChangeRoute(&help)

	case keys.UpdateItemMsg:
		return msg.Cmd(c.Current())

	case keys.ToggleItemsMsg, keys.ToggleItemMsg:
		c.ToggleItems(c.Current())
		cmds = append(cmds, keys.LineDown)

	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
		if msg.String() == "f1" {
			page := filepath.Base(reactea.CurrentRoute())
			reactea.SetCurrentRoute(filepath.Join("help", page))
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
