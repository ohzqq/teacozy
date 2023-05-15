package cmpnt

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/keys"
)

type Filter struct {
	*List
	input  textinput.Model
	Prefix string
	Style  lipgloss.Style
	keyMap keys.KeyMap
}

func NewFilter(props *teacozy.Page) teacozy.PageComponent {
	c := &Filter{
		List:   NewList(props).(*List),
		input:  textinput.New(),
		Prefix: "> ",
		Style:  lipgloss.NewStyle().Foreground(color.Cyan()),
	}
	c.Pager.SetKeyMap(keys.DefaultKeyMap())

	c.List.SetKeyMap(keys.NewKeyMap(keys.Toggle(), keys.Help()))

	return c
}

func (c *Filter) Mount() (reactea.SomeComponent, tea.Cmd) {
	return c, c.input.Focus()
}

func (c *Filter) Update(msg tea.Msg) tea.Cmd {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, k := range c.keyMap.Keys() {
			if key.Matches(msg, k.Binding) {
				cmds = append(cmds, k.TeaCmd)
			}
		}
	}

	if c.input.Focused() {
		c.input, cmd = c.input.Update(msg)
		c.Props().SetInputValue(c.input.Value())
		cmds = append(cmds, cmd)
		//} else {
	}
	cmds = append(cmds, c.List.Update(msg))

	return tea.Batch(cmds...)
}

func (c *Filter) Render(w, h int) string {
	in := c.input.View()
	li := c.Pager.Render(w, h-1)
	return lipgloss.JoinVertical(lipgloss.Left, in, li)
}
