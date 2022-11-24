package key

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Key struct {
	key.Binding
	Cmd MenuFunc
}

type KeyMap struct {
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

func NewKey(k, h string) *Key {
	return &Key{
		Binding: key.NewBinding(
			key.WithKeys(k),
			key.WithHelp(k, h),
		),
	}
}

func NewBinding(k, help string) key.Binding {
	return key.NewBinding(
		key.WithKeys(k),
		key.WithHelp(k, help),
	)
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

func (k KeyMap) FullHelp() *Info {
	return NewInfo().SetData(k)
}

func (k KeyMap) Get(name string) FieldData {
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

func (k KeyMap) Keys() []string {
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

var Keys = KeyMap{
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
	DeselectAll = NewKey("V", "deselect all items")
	EditField   = NewKey("e", "edit meta")
	Enter       = NewKey("enter", "select item")
	ExitScreen  = NewKey("q", "exit screen")
	FullScreen  = NewKey("f", "full screen")
	HelpKey     = NewKey("?", "help")
	InfoKey     = NewKey("i", "view item meta")
	MenuKey     = NewKey("m", "menu")
	PrevScreen  = NewKey("p", "prev screen")
	Quit        = Key{
		Binding: key.NewBinding(
			key.WithKeys("ctrl+c", "esc", "Q"),
			key.WithHelp("ctrl+c/esc", "quit"),
		)}
	SaveAndExit    = NewKey("ctrl+w", "save and exit")
	ToggleAllItems = NewKey("v", "select all")
	SortList       = NewKey("o", "sort")
	ToggleItem     = NewKey("space", "toggle")
	ToggleItemList = NewKey("x", "toggle item list")
)
