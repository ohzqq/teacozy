package body

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/pagy"
	"github.com/ohzqq/teacozy/util"
)

type App struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[teacozy.Props]

	width  int
	height int

	CurrentItem int

	keyMap keys.KeyMap

	*pagy.Paginator
}

type Props struct {
	teacozy.Props
	Width  int
	Height int
}

func New() *App {
	c := &App{
		width:  util.TermWidth(),
		height: util.TermHeight() - 2,
	}

	return c
}

func (c *App) Init(props teacozy.Props) tea.Cmd {
	c.UpdateProps(props)

	//c.Paginator = pagy.New(c.width, c.Props().Items.Len())
	//c.Paginator.SetKeyMap(keys.VimKeyMap())

	return nil
}

func (c *App) Initializer(props teacozy.Props) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		//p := Props{
		//Props: props,
		//}
		return component, component.Init(props)
	}
}

func (c *App) Update(msg tea.Msg) tea.Cmd {
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

func (c App) Name() string {
	return "body"
}

func (c *App) Render(w, h int) string {
	view := teacozy.Renderer(c.Props(), c.width, c.height)

	return view
}

func (c *App) SetKeyMap(km keys.KeyMap) *App {
	c.Paginator.SetKeyMap(km)
	return c
}

func (c *App) AddKey(k *keys.Binding) *App {
	if !c.keyMap.Contains(k) {
		c.keyMap.AddBinds(k)
	} else {
		c.keyMap.Replace(k)
	}
	return c
}

func (m *App) SetCurrent(idx int) {
	m.CurrentItem = idx
}

func (m *App) Current() int {
	return m.CurrentItem
}
