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
	DeSelectAll Key
	EditField   Key
	Enter       Key
	ExitScreen  Key
	FullScreen  Key
	Help        Key
	Info        Key
	Menu        Key
	PrevScreen  Key
	Quit        Key
	SaveAndExit Key
	SelectAll   Key
	SortList    Key
	ToggleItem  Key
}

var Keys = keys{
	DeSelectAll: Key{Bind: DeSelectAll},
	EditField:   Key{Bind: EditField},
	Enter:       Key{Bind: Enter},
	ExitScreen:  Key{Bind: ExitScreen},
	FullScreen:  Key{Bind: FullScreen},
	Help:        Key{Bind: Help},
	Info:        Key{Bind: InfoKey},
	Menu:        Key{Bind: MenuKey},
	PrevScreen:  Key{Bind: PrevScreen},
	Quit:        Key{Bind: Quit},
	SaveAndExit: Key{Bind: SaveAndExit},
	SelectAll:   Key{Bind: SelectAll},
	SortList:    Key{Bind: SortList},
	ToggleItem:  Key{Bind: ToggleItem},
}

var (
	DeSelectAll = key.NewBinding(
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
	InfoKey = key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "view item info"),
	)
	MenuKey = key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "show menu"),
	)
	PrevScreen = key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "prev screen"),
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
