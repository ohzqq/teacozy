package list

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/keys"
)

type List struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	KeyMap keys.KeyMap
}

type Props struct {
	item.Props
	ToggleItems func(...int)
}

func NewList() *List {
	return &List{
		KeyMap: DefaultKeyMap(),
	}
}

func (c *List) Initialize(a *App) {
	a.Routes["list"] = func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		comp := NewList()
		p := Props{
			Props:       a.ItemProps(),
			ToggleItems: a.ToggleItems,
		}
		return comp, comp.Init(p)
	}
}

func (c *List) Init(props Props) tea.Cmd {
	c.UpdateProps(props)
	return nil
}

func (c *List) Update(msg tea.Msg) tea.Cmd {
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

func (c *List) Render(w, h int) string {
	view := item.Renderer(c.Props().Props, w, h)
	return view
	//return c.Props().View()
}

func DefaultKeyMap() keys.KeyMap {
	return keys.KeyMap{
		keys.Toggle().AddKeys(" "),
		keys.Up().WithKeys("k"),
		keys.Down().WithKeys("j"),
		keys.HalfPgUp().WithKeys("K"),
		keys.HalfPgDown().WithKeys("J"),
		keys.Home().WithKeys("g"),
		keys.End().WithKeys("G"),
		keys.Quit().AddKeys("q"),
		keys.New("ctrl+a", "v").
			WithHelp("toggle all").
			Cmd(keys.ToggleAllItems),
		keys.Esc(),
	}
	//return keys.KeyMap{
	//keys.Toggle(),
	//keys.Up(),
	//keys.Down(),
	//keys.HalfPgUp(),
	//keys.HalfPgDown(),
	//keys.Home(),
	//keys.End(),
	//keys.Quit(),
	//keys.New("ctrl+a").
	//WithHelp("toggle all").
	//Cmd(keys.ToggleAllItems),
	//}
}

func VimKeyMap() keys.KeyMap {
	return keys.KeyMap{
		keys.Toggle().AddKeys(" "),
		keys.Up().WithKeys("k"),
		keys.Down().WithKeys("j"),
		keys.HalfPgUp().WithKeys("K"),
		keys.HalfPgDown().WithKeys("J"),
		keys.Home().WithKeys("g"),
		keys.End().WithKeys("G"),
		keys.Quit().AddKeys("q"),
		keys.New("ctrl+a", "v").
			WithHelp("toggle all").
			Cmd(keys.ToggleAllItems),
	}
}
