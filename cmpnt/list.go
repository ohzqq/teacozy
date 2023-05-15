package cmpnt

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
)

type List struct {
	*Pager

	ConfirmChoices bool
	NumSelected    int
	Limit          int
	NoLimit        bool

	keyMap keys.KeyMap
	Style  Style

	help keys.KeyMap
}

type ListProps struct {
	Selected   func() map[int]struct{}
	ToggleItem func(int)
}

func NewList(props *teacozy.Page) teacozy.PageComponent {
	p := New()
	p.Init(props)
	c := &List{
		Limit: 10,
		Style: DefaultStyle(),
		Pager: p,
	}
	c.Pager.ReadOnly = false

	if c.NoLimit {
		c.Limit = props.Items().Len()
	}

	c.keyMap = keys.NewKeyMap(keys.Toggle().AddKeys(" "), keys.Help())

	return c
}

func (c *List) Update(msg tea.Msg) tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case keys.ToggleItemMsg:
		//fmt.Println(c.Limit)
		c.ToggleItem()
		//fmt.Println(c.NumSelected)
		cmds = append(cmds, keys.LineDown)
		//return keys.LineDown
	case tea.KeyMsg:
		for _, k := range c.keyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	cmds = append(cmds, c.Pager.Update(msg))

	return tea.Batch(cmds...)
}

func (c List) ToggleItem() {
	c.ToggleItems(c.Props().Current())
}

func (c *List) ToggleItems(items ...int) {
	for _, idx := range items {
		c.Props().SetCurrent(idx)
		if c.Props().IsSelected(idx) {
			c.Props().DeselectItem(idx)
			c.NumSelected--
		} else if c.NumSelected < c.Limit {
			c.Props().SelectItem(idx)
			c.NumSelected++
		}
	}
}

//func (m List) Chosen() []map[string]string {
//  var chosen []map[string]string
//  if len(m.Selected) > 0 {
//    for k := range m.Selected {
//      l := m.Choices.Label(k)
//      v := m.Choices.String(k)
//      chosen = append(chosen, map[string]string{l: v})
//    }
//  }
//  return chosen
//}

func (m List) KeyMap() keys.KeyMap {
	return m.keyMap
}

func (c *List) Mount() (reactea.SomeComponent, tea.Cmd) {
	return c, nil
}
