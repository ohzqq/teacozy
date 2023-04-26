package list

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/frame"
	"github.com/ohzqq/teacozy/frame/filter"
	"github.com/ohzqq/teacozy/keys"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[teacozy.Props]

	KeyMap keys.KeyMap
}

func New() *Component {
	return &Component{
		KeyMap: DefaultKeyMap(),
	}
}

func (c *Component) Initializer(props teacozy.Props) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		comp := New()
		km := keys.VimKeyMap()
		f := keys.Filter().Cmd(frame.ChangeRoute(filter.New()))
		km.AddBinds(f)
		props.SetKeyMap(km)
		return comp, comp.Init(props)
	}
}

func (c Component) Name() string {
	return "list"
}

func (c *Component) Init(props teacozy.Props) tea.Cmd {
	c.UpdateProps(props)
	return nil
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)
	//var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, k := range c.KeyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	return teacozy.Renderer(c.Props(), w, h)
}

func DefaultKeyMap() keys.KeyMap {
	km := []*keys.Binding{
		keys.New("ctrl+a", "v").
			WithHelp("toggle all").
			Cmd(keys.ToggleAllItems),
		keys.Esc(),
	}
	m := keys.NewKeyMap(km...)
	f := keys.Filter().Cmd(frame.ChangeRoute(filter.New()))
	m.AddBinds(f)
	return m
}
