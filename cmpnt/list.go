package cmpnt

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/keys"
)

type List struct {
	*Pager

	ConfirmChoices bool
	NumSelected    int
	Limit          int
	NoLimit        bool
	ReadOnly       bool

	Choices teacozy.Items
	keyMap  keys.KeyMap
	Style   Style

	help keys.KeyMap
}

type ListProps struct {
	Selected   func() map[int]struct{}
	ToggleItem func(int)
}

func NewList(p *Pager, choices teacozy.Items) *List {
	c := &List{
		Limit:   10,
		Style:   DefaultStyle(),
		Choices: choices,
	}

	if c.NoLimit {
		c.Limit = c.Choices.Len()
	}

	c.AddKey(keys.Toggle().AddKeys(" "))

	c.SetKeyMap(keys.VimKeyMap())

	c.AddKey(keys.Help())

	return c
}

func (c *List) Update(msg tea.Msg) tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, k := range c.keyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	return tea.Batch(cmds...)
}

func (c List) ToggleItem() {
	//c.ToggleItems(c.Current())
}

//func (c *List) ToggleItems(items ...int) {
//  for _, idx := range items {
//    c.Props().SetCurrent(idx)
//    if _, ok := c.Selected[idx]; ok {
//      delete(c.Selected, idx)
//      c.NumSelected--
//    } else if c.NumSelected < c.Limit {
//      c.Selected[idx] = struct{}{}
//      c.NumSelected++
//    }
//  }
//}

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
