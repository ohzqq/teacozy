package list

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type List struct {
	Items         Items
	IsPrompt      bool
	IsMultiSelect bool
}

func UpdateList(m *Model, msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.IsMulti() {
			switch {
			case key.Matches(msg, m.Keys.Enter):
				if m.ShowSelectedOnly {
					cmds = append(cmds, ReturnSelectionsCmd())
				}
				m.ShowSelectedOnly = true
				cmds = append(cmds, UpdateVisibleItemsCmd("selected"))
			case key.Matches(msg, m.Keys.SelectAll):
				ToggleAllItemsCmd(m)
				cmds = append(cmds, UpdateVisibleItemsCmd("all"))
			}
		} else {
			switch {
			case key.Matches(msg, m.Keys.Enter):
				cur := m.List.SelectedItem().(Item)
				m.SetItem(m.List.Index(), cur.ToggleSelected())
				cmds = append(cmds, ReturnSelectionsCmd())
			}
		}

		switch {
		case key.Matches(msg, m.Keys.ExitScreen):
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, m.Keys.Prev):
			m.ShowSelectedOnly = false
			cmds = append(cmds, UpdateVisibleItemsCmd("all"))
		}
	}
	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}
