package key

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy"
)

type Key struct {
	key teacozy.Key
	//Cmd MenuFunc
}

func Matches(msg tea.KeyMsg, bind ...key.Binding) bool {
	return key.Matches(msg, bind...)
}

func NewKey(k, h string) *Key {
	return &Key{
		key: teacozy.Key{
			Binding: key.NewBinding(
				key.WithKeys(k),
				key.WithHelp(k, h),
			),
		},
	}
}

func NewBinding(k, help string) key.Binding {
	return key.NewBinding(
		key.WithKeys(k),
		key.WithHelp(k, help),
	)
}

func (k Key) Binding() key.Binding {
	return k.key.Binding
}

func (k Key) Cmd() teacozy.MenuFunc {
	return k.key.Cmd
}

func (k *Key) SetCmd(cmd teacozy.MenuFunc) *Key {
	k.key.Cmd = cmd
	return k
}

func (k Key) Matches(msg tea.KeyMsg) bool {
	return key.Matches(msg, k.key.Binding)
}

func (i Key) Key() string {
	return i.key.Binding.Help().Key
}

func (i Key) Value() string {
	return i.key.Binding.Help().Desc
}

func (i Key) Set(v string) {}

func (i Key) String() string {
	return i.key.Binding.Help().Key + ": " + i.key.Binding.Help().Desc
}

type KeyMap struct {
	keys []*Key
}

func NewKeyMap() KeyMap {
	return KeyMap{}
}

func (k KeyMap) All() []*Key {
	return k.keys
}

func (k *KeyMap) New(key, help string) {
	bind := NewKey(key, help)
	k.Add(bind)
}

func (k *KeyMap) Add(key *Key) {
	k.keys = append(k.keys, key)
}

//func (k KeyMap) FullHelp() *info.Info {
//  return info.NewInfo().SetData(k)
//}

func (k KeyMap) Get(name string) teacozy.FieldData {
	var key *Key
	for _, bind := range k.keys {
		if bind.Key() == name {
			key = bind
		}
	}
	return key
}

func (k KeyMap) Keys() []string {
	var keys []string
	for _, bind := range k.keys {
		keys = append(keys, bind.Key())
	}
	return keys
}

var Map = teacozy.KeysMap{
	DeselectAll:    teacozy.Key{Binding: DeselectAll},
	EditField:      teacozy.Key{Binding: EditField},
	Enter:          teacozy.Key{Binding: Enter},
	ExitScreen:     teacozy.Key{Binding: ExitScreen},
	FullScreen:     teacozy.Key{Binding: FullScreen},
	Help:           teacozy.Key{Binding: HelpKey},
	Info:           teacozy.Key{Binding: InfoKey},
	Menu:           teacozy.Key{Binding: MenuKey},
	PrevScreen:     teacozy.Key{Binding: PrevScreen},
	Quit:           teacozy.Key{Binding: Quit},
	SaveAndExit:    teacozy.Key{Binding: SaveAndExit},
	ToggleAllItems: teacozy.Key{Binding: ToggleAllItems},
	SortList:       teacozy.Key{Binding: SortList},
	ToggleItem:     teacozy.Key{Binding: ToggleItem},
	ToggleItemList: teacozy.Key{Binding: ToggleItemList},
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
