package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/bubbles/key"
	"github.com/ohzqq/teacozy/list"
)

type Model struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	router      reactea.Component[router.Props]
	items       *list.Items
	opts        []list.Option
	currentItem *list.Item
	KeyMap      KeyMap
}

type KeyMap struct {
	Edit key.Binding
	Prev key.Binding
}

func New(items *list.Items, opts ...list.Option) *Model {
	m := &Model{
		router: router.New(),
		items:  items,
		opts:   opts,
		KeyMap: DefaultKeyMap(),
	}
	m.SetCurrentItem(items.Get(0))
	return m
}

func (m *Model) Init(reactea.NoProps) tea.Cmd {
	return m.router.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			c := NewList()
			return c, c.Init(Props{
				Items:          m.items,
				Opts:           m.opts,
				SetCurrentItem: m.SetCurrentItem,
			})
		},
		"edit": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			c := NewList()
			opts := []list.Option{list.Editable(true)}
			opts = append(opts, m.opts...)
			return c, c.Init(Props{
				Items:          m.items,
				Opts:           opts,
				SetCurrentItem: m.SetCurrentItem,
			})
		},
	})
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}
		switch {
		case key.Matches(msg, m.KeyMap.Edit):
			reactea.SetCurrentRoute("edit")
		}
		if msg.String() == "p" {
			switch reactea.CurrentRoute() {
			case "edit":
				m.items.SetEditable(false)
				reactea.SetCurrentRoute("default")
			}
		}
	}

	return m.router.Update(msg)
}

func (m *Model) Render(w, h int) string {
	view := m.router.Render(w, h)
	return view
}

func (m *Model) SetCurrentItem(li *list.Item) {
	m.currentItem = li
}

func DefaultKeyMap() KeyMap {
	km := KeyMap{
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit list"),
		),
		Prev: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "prev page"),
		),
	}
	//km.Edit.SetEnabled(false)
	return km
}
