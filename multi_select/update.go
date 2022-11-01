package multi

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Enter):
			if m.ShowSelectedOnly {
				cmds = append(cmds, ReturnSelectionsCmd())
			}

			m.ShowSelectedOnly = true
			cmds = append(cmds, UpdateDisplayedItemsCmd("selected"))
		case key.Matches(msg, m.Keys.SelectAll):
			ToggleAllItemsCmd(m)
			cmds = append(cmds, UpdateDisplayedItemsCmd("all"))
		}
	}

	return m, tea.Batch(cmds...)
}
