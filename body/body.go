package body

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/util"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[teacozy.Props]

	width  int
	height int

	CurrentItem int

	keyMap keys.KeyMap
}

type Props struct {
	teacozy.Props
	Width  int
	Height int
}

func New() *Component {
	c := &Component{
		width:  util.TermWidth(),
		height: util.TermHeight() - 2,
	}

	return c
}

func (c *Component) Init(props teacozy.Props) tea.Cmd {
	c.UpdateProps(props)
	return nil
}

func (c *Component) Initializer(props teacozy.Props) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		return component, component.Init(props)
	}
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// ctrl+c support
		if msg.String() == "ctrl+c" {
			return reactea.Destroy
		}

		for _, k := range c.keyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	p := c.Props().Paginator

	p, cmd = p.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (c Component) Name() string {
	return "body"
}

func (c *Component) Render(w, h int) string {
	view := teacozy.Renderer(c.Props(), c.width, c.height)
	return view
}

func (c *Component) SetKeyMap(km keys.KeyMap) *Component {
	c.Paginator.SetKeyMap(km)
	return c
}

func (c *Component) AddKey(k *keys.Binding) *Component {
	if !c.keyMap.Contains(k) {
		c.keyMap.AddBinds(k)
	} else {
		c.keyMap.Replace(k)
	}
	return c
}

func (m *Component) SetCurrent(idx int) {
	m.CurrentItem = idx
}

func (m *Component) Current() int {
	return m.CurrentItem
}
