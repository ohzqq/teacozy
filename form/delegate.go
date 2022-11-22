package form

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func itemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string
		if i, ok := m.SelectedItem().(*Field); ok {
			title = i.Title()
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, EditField):
				return EditFormItemCmd(curItem)
			}
		}
		return nil
	}

	return d
}

var (
	DeselectAll = key.NewBinding(
		key.WithKeys("V"),
		key.WithHelp("V", "deselect all items"),
	)
	EditField = key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit meta"),
	)
	Enter = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select item"),
	)
	ExitScreen = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "exit screen"),
	)
	FullScreen = key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "full screen"),
	)
	HelpKey = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	)
	InfoKey = key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "view item meta"),
	)
	MenuKey = key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "menu"),
	)
	PrevScreen = key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "prev screen"),
	)
	Quit = key.NewBinding(
		key.WithKeys("ctrl+c", "esc", "Q"),
		key.WithHelp("ctrl+c/esc", "quit"),
	)
	SaveAndExit = key.NewBinding(
		key.WithKeys("ctrl+w"),
		key.WithHelp("ctrl+w", "save and exit"),
	)
	ToggleAllItems = key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "select all"),
	)
	SortList = key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "sort"),
	)
	ToggleItem = key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle"),
	)
	ToggleItemList = key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "toggle item list"),
	)
)
