package cmpnt

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/keys"
)

type Help struct {
	*Pager

	keyMap keys.KeyMap
}

func NewHelp(props Props) Page {
	p := New()
	p.Init(props)
	help := &Help{
		keyMap: keys.NewKeyMap(keys.Esc()),
		Pager:  p,
	}
	return help
}

func (c *Help) Component() reactea.SomeComponent {
	return c
}

func (c Help) KeyMap() keys.KeyMap {
	return c.keyMap
}

func (c *Help) Update(msg tea.Msg) tea.Cmd {
	reactea.AfterUpdate(c)
	var cmds []tea.Cmd

	switch msg := msg.(type) {
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
