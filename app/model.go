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
}

type List struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	list *list.Model
}

type Props struct {
	Items *list.Items
	Opts  []list.Option
}

func New(items *list.Items, opts ...list.Option) *Model {
	return &Model{
		router: router.New(),
		items:  items,
		opts:   opts,
	}
}

func NewList() *List {
	return &List{}
}

func (l *List) Init(props Props) tea.Cmd {
	l.list = list.New(props.Items, props.Opts...)
	l.UpdateProps(props)
}

func (l *List) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	l.list, cmd = l.list.Update(msg)
	return cmd
}

func (l *List) Render(w, h int) string {
	return l.list.View()
}

func (m *Model) Init(reactea.NoProps) tea.Cmd {
	return m.router.Init(map[string]router.RouteInitializer{
		"default": func(router.Params) (reactea.SomeComponent, tea.Cmd) {
			c := NewList()
			return c, c.Init(Props{Items: m.items, Opts: m.opts})
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

	return c.router.Update(msg)
}

func (m *Model) Render(w, h int) string {
	return c.router.Render(w, h)
}
