package choose

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
)

func (m *Choose) ReturnSelections() tea.Cmd {
	//if m.Props().Limit == 1 {
	//  return message.ToggleItem()
	//}
	//if m.Props().NumSelected == 0 {
	//  m.quitting = true
	//  return message.ToggleItem()
	//}
	return message.ReturnSels(m.Props().Limit, m.Props().NumSelected)
}

func (m *Choose) quit() tea.Cmd {
	m.quitting = true
	return message.ReturnSelections()
}

func (m *Choose) KeyMap() keys.KeyMap {
	var keys = keys.KeyMap{
		keys.ShowHelp(),
		keys.Quit().
			WithKeys("ctrl+c", "q").
			Cmd(m.quit()),
		keys.ToggleItem().WithKeys("tab", " "),
		keys.NewBinding("e").
			WithHelp("edit field").
			Cmd(message.StartEditing()),
		keys.NewBinding("enter").
			WithHelp("return selections").
			Cmd(m.ReturnSelections()),
		keys.NewBinding("/").
			WithHelp("filter list").
			Cmd(message.StartFiltering()),
	}
	return keys
}

// key.NewBinding(
//
//	key.WithKeys("V"),
//	key.WithHelp("V", "deselect all"),
//
// ),
// key.NewBinding(
//
//	key.WithKeys("v"),
//	key.WithHelp("v", "select all"),
//
// ),
