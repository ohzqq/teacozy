package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/list"
)

type Model struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[reactea.NoProps]

	router reactea.Component[router.Props]
	items  *list.Items
	opts   []list.Option
	List   *list.Model
}

type List struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	list *list.Model
}

type Props struct {
	SetItems func(*list.Items)
}

func New(items *list.Items, opts ...list.Option) *Model {
	return &Model{
		router: router.New(),
		items:  items,
		opts:   opts,
	}
}

func NewList(items *list.Items, opts []list.Option) *List {
	return &List{
		list: list.New(items, opts...),
	}
}

func (l *List) Init(props Props) tea.Cmd {
	println(l.list.Len())
	l.UpdateProps(props)
	return nil
}

func (l *List) Update(msg tea.Msg) tea.Cmd {
	m, cmd := l.list.Update(msg)
	l.list = m.(*list.Model)
	//l.Props().SetItems(l.list.Items)
	return cmd
}

func (l *List) Render(w, h int) string {
	l.list.SetSize(w, h)
	return l.list.View()
}

func (m *Model) Init(reactea.NoProps) tea.Cmd {
	return m.router.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			c := NewList(m.items, m.opts)
			return c, c.Init(Props{SetItems: m.SetItems})
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
	}

	return m.router.Update(msg)
}

func (m *Model) Render(w, h int) string {
	return m.router.Render(w, h)
}

func (m *Model) SetItems(items *list.Items) {
	m.items = items
}
