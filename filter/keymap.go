package filter

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/keymap"
)

func FilterKeyMap(m *Model) keymap.KeyMap {
	km := keymap.KeyMap{
		keymap.NewBinding(
			keymap.WithKeys("esc"),
			keymap.WithHelp("esc", "stop filtering"),
			keymap.WithCmd(StopFilteringCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("enter"),
			keymap.WithHelp("enter", "return selections"),
			keymap.WithCmd(StopFilteringCmd(m)),
		),
	}
	return km
}

func GlobalKeyMap(m *Model) keymap.KeyMap {
	return keymap.KeyMap{
		keymap.NewBinding(
			keymap.WithKeys("down"),
			keymap.WithHelp("down", "move cursor down"),
			keymap.WithCmd(DownCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("up"),
			keymap.WithHelp("up", "move cursor up"),
			keymap.WithCmd(UpCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("ctrl+c"),
			keymap.WithHelp("ctrl+c", "quit"),
			keymap.WithCmd(tea.Quit),
		),
		keymap.NewBinding(
			keymap.WithKeys("tab"),
			keymap.WithHelp("tab", "select item"),
			keymap.WithCmd(SelectItemCmd(m)),
		),
	}
}

func ListKeyMap(m *Model) keymap.KeyMap {
	km := keymap.KeyMap{
		keymap.NewBinding(
			keymap.WithKeys("right", "l"),
			keymap.WithHelp("right/l", "next page"),
			keymap.WithCmd(NextPageCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("left", "h"),
			keymap.WithHelp("left/h", "prev page"),
			keymap.WithCmd(PrevPageCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("V"),
			keymap.WithHelp("V", "deselect all"),
			keymap.WithCmd(DeselectAllItemsCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("v"),
			keymap.WithHelp("v", "select all"),
			keymap.WithCmd(SelectAllItemsCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys(" "),
			keymap.WithHelp("space", "select item"),
			keymap.WithCmd(SelectItemCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("j"),
			keymap.WithHelp("j", "move cursor down"),
			keymap.WithCmd(DownCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("k"),
			keymap.WithHelp("k", "move cursor up"),
			keymap.WithCmd(UpCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("esc", "q"),
			keymap.WithHelp("esc/q", "quit"),
			keymap.WithCmd(tea.Quit),
		),
		keymap.NewBinding(
			keymap.WithKeys("enter"),
			keymap.WithHelp("enter", "return selections"),
			keymap.WithCmd(EnterCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("/"),
			keymap.WithHelp("/", "filter items"),
			keymap.WithCmd(FilterItemsCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("G"),
			keymap.WithHelp("G", "last item"),
			keymap.WithCmd(BottomCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("g"),
			keymap.WithHelp("g", "first item"),
			keymap.WithCmd(TopCmd(m)),
		),
	}
	return km
}
