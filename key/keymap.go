package key

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/ohzqq/teacozy"
)

type KeyMap struct {
	keys []*Key
}

func NewKeyMap() KeyMap {
	return KeyMap{}
}

func (k KeyMap) All() []*Key {
	return k.keys
}

func (k *KeyMap) NewKey(key, help string) {
	bind := NewKey(key, help)
	k.Add(bind)
}

func (k *KeyMap) Add(key *Key) {
	k.keys = append(k.keys, key)
}

func (k *KeyMap) AddBind(kb key.Binding) {
	bind := &Key{key: kb}
	k.Add(bind)
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

func (k KeyMap) ShortHelp() []key.Binding {
	var keys []key.Binding
	for _, bind := range k.keys {
		keys = append(keys, bind.Binding())
	}
	return keys
}

func (k KeyMap) Get(name string) teacozy.Field {
	var key *Key
	for _, bind := range k.keys {
		if bind.Name() == name {
			key = bind
		}
	}
	return key
}

func (k KeyMap) GetKey(name string) *Key {
	var key *Key
	for _, bind := range k.keys {
		if bind.Name() == name {
			key = bind
		}
	}
	return key
}

func (k KeyMap) Keys() []string {
	var keys []string
	for _, bind := range k.keys {
		keys = append(keys, bind.Name())
	}
	return keys
}

//var Map = teacozy.KeysMap{
//  DeselectAll:    teacozy.Key{Binding: DeselectAll},
//  EditField:      teacozy.Key{Binding: EditField},
//  Enter:          teacozy.Key{Binding: Enter},
//  ExitScreen:     teacozy.Key{Binding: ExitScreen},
//  FullScreen:     teacozy.Key{Binding: FullScreen},
//  Help:           teacozy.Key{Binding: HelpKey},
//  Info:           teacozy.Key{Binding: InfoKey},
//  Menu:           teacozy.Key{Binding: MenuKey},
//  PrevScreen:     teacozy.Key{Binding: PrevScreen},
//  Quit:           teacozy.Key{Binding: Quit},
//  SaveAndExit:    teacozy.Key{Binding: SaveAndExit},
//  ToggleAllItems: teacozy.Key{Binding: ToggleAllItems},
//  SortList:       teacozy.Key{Binding: SortList},
//  ToggleItem:     teacozy.Key{Binding: ToggleItem},
//  ToggleItemList: teacozy.Key{Binding: ToggleItemList},
//}

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
