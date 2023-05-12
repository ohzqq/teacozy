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

func NewHelp(maps ...keys.KeyMap) *Help {
	var km keys.KeyMap
	for _, m := range maps {
		km.AddBinds(m.Keys()...)
	}

	help := &Help{
		keyMap: keys.NewKeyMap(keys.Esc()),
	}

	help.Pager = New(
		WithMap(km.Map()),
		ReadOnly(),
	)

	return help
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
