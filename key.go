package teacozy

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Key struct {
	key.Binding
	Cmd MenuFunc
}

func NewKey(k, h string) *Key {
	return &Key{
		Binding: key.NewBinding(
			key.WithKeys(k),
			key.WithHelp(k, h),
		),
	}
}

func (k *Key) SetCmd(cmd MenuFunc) *Key {
	k.Cmd = cmd
	return k
}

func (k Key) Matches(msg tea.KeyMsg) bool {
	return key.Matches(msg, k.Binding)
}

func (i Key) Key() string {
	return i.Binding.Help().Key
}

func (i Key) Value() string {
	return i.Binding.Help().Desc
}

func (i Key) Set(v string) {}

func (i Key) String() string {
	return i.Binding.Help().Key + ": " + i.Binding.Help().Desc
}

func NewKeyBind(k, help string) key.Binding {
	return key.NewBinding(
		key.WithKeys(k),
		key.WithHelp(k, help),
	)
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

type KeysMap struct {
	DeselectAll    Key
	EditField      Key
	Enter          Key
	ExitScreen     Key
	FullScreen     Key
	Help           Key
	Info           Key
	Menu           Key
	PrevScreen     Key
	Quit           Key
	SaveAndExit    Key
	ToggleAllItems Key
	SortList       Key
	ToggleItem     Key
	ToggleItemList Key
}

type KeyMap map[string]key.Binding

func (k KeysMap) FullHelp() *Info {
	return NewInfo().SetData(k)
}

func (k KeysMap) Get(name string) FieldData {
	var key Key
	switch name {
	case "Deselect All Items":
		key = k.DeselectAll
	case "Edit Field":
		key = k.EditField
	case "Enter":
		key = k.Enter
	case "Exit Screen":
		key = k.ExitScreen
	case "Full Screen":
		key = k.FullScreen
	case "Help":
		key = k.Help
	case "Item Meta":
		key = k.Info
	case "Main Menu":
		key = k.Menu
	case "Prev Screen":
		key = k.PrevScreen
	case "Quit":
		key = k.Quit
	case "Save And Exit":
		key = k.SaveAndExit
	case "Select All":
		key = k.ToggleAllItems
	case "Sort List":
		key = k.SortList
	case "Toggle Item":
		key = k.ToggleItem
	case "Toggle Item List":
		key = k.ToggleItemList
	}
	return key
}

func (k KeysMap) Keys() []string {
	return []string{
		"Toggle Item",
		"Quit",
		"Save And Exit",
		"Edit Field",
		"Enter",
		"Full Screen",
		"Item Meta",
		"Main Menu",
		"Sort List",
		"Prev Screen",
		"Exit Screen",
		"Deselect All Items",
		"Toggle All Items",
		"Toggle Item List",
		"Help",
	}
}

var Keys = KeysMap{
	DeselectAll:    Key{Binding: DeselectAll},
	EditField:      Key{Binding: EditField},
	Enter:          Key{Binding: Enter},
	ExitScreen:     Key{Binding: ExitScreen},
	FullScreen:     Key{Binding: FullScreen},
	Help:           Key{Binding: HelpKey},
	Info:           Key{Binding: InfoKey},
	Menu:           Key{Binding: MenuKey},
	PrevScreen:     Key{Binding: PrevScreen},
	Quit:           Key{Binding: Quit},
	SaveAndExit:    Key{Binding: SaveAndExit},
	ToggleAllItems: Key{Binding: ToggleAllItems},
	SortList:       Key{Binding: SortList},
	ToggleItem:     Key{Binding: ToggleItem},
	ToggleItemList: Key{Binding: ToggleItemList},
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
