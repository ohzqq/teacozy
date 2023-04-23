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
	reactea.BasicPropfulComponent[Props]

	KeyMap  keys.KeyMap
	current int
}

type Props struct {
	teacozy.Props
	ToggleItem func()
}

func New() *Component {
	return &Component{
		KeyMap: DefaultKeyMap(),
	}
}

func (c *Component) Initialize(a *frame.App) {
	a.Routes["list"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		comp := New()
		p := Props{
			Props:      a.ItemProps(),
			ToggleItem: a.ToggleItem,
		}
		a.SetKeyMap(frame.DefaultKeyMap())
		return comp, comp.Init(p)
	}
}

func (c Component) Name() string {
	return "list"
}

func (c *Component) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	return nil
}

func (c *Component) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)
	//var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "/" {
			return frame.ChangeRoute(filter.New())
		}
		for _, k := range c.KeyMap {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	props := c.Props().Props
	return teacozy.Renderer(props, w, h)
}

func (c *Component) setCurrent(i int) {
	c.current = i
}

func DefaultKeyMap() keys.KeyMap {
	km := keys.KeyMap{
		keys.Toggle().AddKeys(" "),
		keys.New("ctrl+a", "v").
			WithHelp("toggle all").
			Cmd(keys.ToggleAllItems),
		keys.Esc(),
	}
	return km
}
