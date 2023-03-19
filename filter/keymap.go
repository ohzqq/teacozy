package filter

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/keys"
)

func FilterKeyMap(m *Model) keys.KeyMap {
	km := keys.KeyMap{
		keys.NewBinding(
			keys.WithKeys("esc"),
			keys.WithHelp("esc", "stop filtering"),
			keys.WithCmd(StopFilteringCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("enter"),
			keys.WithHelp("enter", "return selections"),
			keys.WithCmd(StopFilteringCmd(m)),
		),
	}
	return km
}

func GlobalKeyMap(m *Model) keys.KeyMap {
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
			keys.WithCmd(tea.Quit),
		),
		keys.NewBinding(
			keys.WithKeys("tab"),
			keys.WithHelp("tab", "select item"),
			keys.WithCmd(SelectItemCmd(m)),
		),
	}
}

func ListKeyMap(m *Model) keys.KeyMap {
	km := keys.KeyMap{
		keys.NewBinding(
			keys.WithKeys("right", "l"),
			keys.WithHelp("right/l", "next page"),
			keys.WithCmd(NextPageCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("left", "h"),
			keys.WithHelp("left/h", "prev page"),
			keys.WithCmd(PrevPageCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("V"),
			keys.WithHelp("V", "deselect all"),
			keys.WithCmd(DeselectAllItemsCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("v"),
			keys.WithHelp("v", "select all"),
			keys.WithCmd(SelectAllItemsCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys(" "),
			keys.WithHelp("space", "select item"),
			keys.WithCmd(SelectItemCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("j"),
			keys.WithHelp("j", "move cursor down"),
			keys.WithCmd(DownCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("k"),
			keys.WithHelp("k", "move cursor up"),
			keys.WithCmd(UpCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("esc", "q"),
			keys.WithHelp("esc/q", "quit"),
			keys.WithCmd(tea.Quit),
		),
		keys.NewBinding(
			keys.WithKeys("enter"),
			keys.WithHelp("enter", "return selections"),
			keys.WithCmd(EnterCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("/"),
			keys.WithHelp("/", "filter items"),
			keys.WithCmd(FilterItemsCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("G"),
			keys.WithHelp("G", "last item"),
			keys.WithCmd(BottomCmd(m)),
		),
		keys.NewBinding(
			keys.WithKeys("g"),
			keys.WithHelp("g", "first item"),
			keys.WithCmd(TopCmd(m)),
		),
	}
	return km
}
