package list

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/frame"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/keys"
)

type Component struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	KeyMap keys.KeyMap
}

type Props struct {
	item.Props
	ToggleItems func(...int)
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
			Props:       a.ItemProps(),
			ToggleItems: a.ToggleItems,
		}
		return comp, comp.Init(p)
	}
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
	case keys.ToggleItemMsg:
		c.Props().ToggleItems(c.Props().Cursor())
		cmds = append(cmds, keys.LineDown)
	case tea.KeyMsg:
		for _, k := range c.KeyMap {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	return tea.Batch(cmds...)
}

func (c *Component) Render(w, h int) string {
	view := item.Renderer(c.Props().Props, w, h)
	return view
	//return c.Props().View()
}

func DefaultKeyMap() keys.KeyMap {
	return keys.KeyMap{
		keys.Toggle().AddKeys(" "),
		keys.New("ctrl+a", "v").
			WithHelp("toggle all").
			Cmd(keys.ToggleAllItems),
		keys.Esc(),
	}
}
