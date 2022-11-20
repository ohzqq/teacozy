package teacozy

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Key struct {
	Bind key.Binding
	Cmd  MenuFunc
}

func NewKey(k, h string) *Key {
	return &Key{
		Bind: key.NewBinding(
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
	return key.Matches(msg, k.Bind)
}

func (i Key) Key() string {
	return i.Bind.Help().Key
}

func (i Key) Value() string {
	return i.Bind.Help().Desc
}

func (i Key) Set(v string) {}

func (i Key) String() string {
	return i.Bind.Help().Key + ": " + i.Bind.Help().Desc
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

type keys struct {
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
}

func (k keys) FullHelp() *Info {
	return NewInfo().SetData(k)
}

func (k keys) Get(name string) FieldData {
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
	}
	return key
}

func (k keys) Keys() []string {
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
		"Help",
	}
}

var Keys = keys{
	DeselectAll:    Key{Bind: DeselectAll},
	EditField:      Key{Bind: EditField},
	Enter:          Key{Bind: Enter},
	ExitScreen:     Key{Bind: ExitScreen},
	FullScreen:     Key{Bind: FullScreen},
	Help:           Key{Bind: HelpKey},
	Info:           Key{Bind: InfoKey},
	Menu:           Key{Bind: MenuKey},
	PrevScreen:     Key{Bind: PrevScreen},
	Quit:           Key{Bind: Quit},
	SaveAndExit:    Key{Bind: SaveAndExit},
	ToggleAllItems: Key{Bind: ToggleAllItems},
	SortList:       Key{Bind: SortList},
	ToggleItem:     Key{Bind: ToggleItem},
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
)
