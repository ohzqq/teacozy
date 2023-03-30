package filter

import (
	"github.com/ohzqq/teacozy/keys"
)

func FilterKeyMap(m *Filter) keys.KeyMap {
	km := keys.KeyMap{
		keys.NewBinding(
			keys.WithKeys("esc"),
			keys.WithHelp("esc", "stop filtering"),
			keys.WithCmd(StopFilteringCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("enter"),
			keys.WithHelp("enter", "return selections"),
			keys.WithCmd(ReturnSelectionsCmd(m)),
		),
	}
	return km
}

func GlobalKeyMap(m *Filter) keys.KeyMap {
	return keys.KeyMap{
		keys.NewBinding(
			keys.WithKeys("down"),
			keys.WithHelp("down", "move cursor down"),
			keys.WithCmd(DownCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("up"),
			keys.WithHelp("up", "move cursor up"),
			keys.WithCmd(UpCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("ctrl+c"),
			keys.WithHelp("ctrl+c", "quit"),
			keys.WithCmd(QuitCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("tab"),
			keys.WithHelp("tab", "select item"),
			keys.WithCmd(SelectItemCmd(m)),
		),
	}
}
