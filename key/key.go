package key

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

const tab = ` `

type KeyMap struct {
	ToggleItem  key.Binding
	SelectAll   key.Binding
	DeselectAll key.Binding
	Unfocus     key.Binding
	Switch      key.Binding
	Enter       key.Binding
	ExitScreen  key.Binding
	Prev        key.Binding
	list.KeyMap
}

func DefaultKeys() KeyMap {
	return KeyMap{
		KeyMap: ListKeyMap(),
		Unfocus: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "exit view"),
		),
		ToggleItem:  ToggleItem,
		SelectAll:   SelectAll,
		DeselectAll: DeselectAll,
		Enter:       Enter,
		ExitScreen:  ExitScreen,
		Prev:        PrevMenu,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	var keys [][]key.Binding
	first := []key.Binding{}
	keys = append(keys, first)
	second := []key.Binding{}
	keys = append(keys, second)
	third := []key.Binding{
		k.ToggleItem,
		k.DeselectAll,
		k.SelectAll,
	}
	keys = append(keys, third)
	return keys
}

func (s KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func ListKeyMap() list.KeyMap {
	km := list.DefaultKeyMap()
	km.NextPage = key.NewBinding(
		key.WithKeys("right", "l", "pgdown"),
		key.WithHelp("l/pgdn", "next page"),
	)
	km.Quit = key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c", "quit"),
	)
	return km
}

var (
	Quit = key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c", "quit"),
	)
	SaveAndExit = key.NewBinding(
		key.WithKeys("ctrl+w"),
		key.WithHelp("ctrl+w", "save and exit"),
	)
	Enter = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select item"),
	)
	FullHelp = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "full help"),
	)
	ToggleItem = key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "select item"),
	)
	PrevMenu = key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "prev menu"),
	)
	ExitScreen = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "exit screen"),
	)
	DeselectAll = key.NewBinding(
		key.WithKeys("V"),
		key.WithHelp("V", "deselect all"),
	)
	SelectAll = key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "select all"),
	)
	ChangeLibrary = key.NewBinding(
		key.WithKeys("L"),
		key.WithHelp("L", "change library"),
	)
	FullScreen = key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "full screen"),
	)
	SortBy = key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "sort options"),
	)
	EditField = key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit meta"),
	)
	Info = key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "view item info"),
	)
	CategoryList = key.NewBinding(
		key.WithKeys("c", "tab"),
		key.WithHelp("c", "Browse Categories"),
	)
)
