package key

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/data"
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

func (k KeyMap) Matches(msg tea.KeyMsg) bool {
	return k.GetKey(msg.String()).Matches(msg)
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

func (k KeyMap) Get(name string) data.Field {
	var key *Key
	for _, bind := range k.keys {
		if bind.Key() == name {
			key = bind
		}
	}
	return key
}

func (k KeyMap) GetKey(name string) *Key {
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

var (
	ToggleItem = key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "toggle"),
	)
	HelpKey = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	)
	Quit = key.NewBinding(
		key.WithKeys("ctrl+c", "esc", "Q"),
		key.WithHelp("ctrl+c/esc", "quit"),
	)
	SaveAndExit = key.NewBinding(
		key.WithKeys("ctrl+w"),
		key.WithHelp("ctrl+w", "save and exit"),
	)
	EditField = key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit meta"),
	)
	Enter = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select item"),
	)
	FullScreen = key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "full screen"),
	)
	InfoKey = key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "view item meta"),
	)
	MenuKey = key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "menu"),
	)
	SortList = key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "sort"),
	)
	PrevScreen = key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "prev screen"),
	)
	ExitScreen = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "exit screen"),
	)
	ToggleAllItems = key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "select all"),
	)
	UnToggleAllItems = key.NewBinding(
		key.WithKeys("V"),
		key.WithHelp("V", "deselect all items"),
	)
	ToggleItemList = key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "toggle item list"),
	)
	Up = key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("???/k", "up"),
	)
	Down = key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("???/j", "down"),
	)
)
