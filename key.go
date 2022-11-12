package teacozy

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

type keys struct {
	DeselectAll key.Binding
	EditField   key.Binding
	Enter       key.Binding
	ExitScreen  key.Binding
	FullScreen  key.Binding
	Help        key.Binding
	Info        key.Binding
	PrevScreen  key.Binding
	Quit        key.Binding
	SaveAndExit key.Binding
	SelectAll   key.Binding
	SortList    key.Binding
	ToggleItem  key.Binding
}

var Keys = keys{
	DeselectAll: DeleselectAll,
	EditField:   EditField,
	Enter:       Enter,
	ExitScreen:  ExitScreen,
	FullScreen:  FullScreen,
	Help:        Help,
	Info:        Info,
	PrevScreen:  PrevScreen,
	Quit:        Quit,
	SaveAndExit: SaveAndExit,
	SelectAll:   SelectAll,
	SortList:    SortList,
	ToggleItem:  ToggleItem,
}

var (
	DeselectAll = key.NewBinding(
		key.WithKeys("V"),
		key.WithHelp("V", "deselect all"),
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
	Help = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "full help"),
	)
	Info = key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "view item info"),
	)
	PrevScreen = key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "prev menu"),
	)
	Quit = key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c", "quit"),
	)
	SaveAndExit = key.NewBinding(
		key.WithKeys("ctrl+w"),
		key.WithHelp("ctrl+w", "save and exit"),
	)
	SelectAll = key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "select all"),
	)
	SortList = key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "sort options"),
	)
	ToggleItem = key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "select item"),
	)
)
