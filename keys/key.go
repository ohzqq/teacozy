package keys

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/message"
)

type KeyMap []*Binding

type Binding struct {
	key.Binding
	help   string
	keys   []key.Binding
	TeaCmd tea.Cmd
}

func NewBinding(keys ...string) *Binding {
	k := Binding{
		Binding: key.NewBinding(),
	}
	k.WithKeys(keys...)
	return &k
}

func (k *Binding) Cmd(cmd tea.Cmd) *Binding {
	k.TeaCmd = cmd
	return k
}

func (k *Binding) WithKeys(keys ...string) *Binding {
	k.Binding.SetKeys(keys...)
	k.WithHelp(k.help)
	return k
}

func (k *Binding) AddKeys(keys ...string) *Binding {
	keys = append(keys, k.Binding.Keys()...)
	k.Binding.SetKeys(keys...)
	k.WithHelp(k.help)
	return k
}

func (k *Binding) WithHelp(h string) *Binding {
	k.help = h
	k.Binding.SetHelp(strings.Join(k.Keys(), "/"), h)
	return k
}

func (km KeyMap) Map() []map[string]string {
	c := make([]map[string]string, len(km))
	for i, k := range km {
		c[i] = map[string]string{k.Help().Key: k.Help().Desc}
	}
	return c
}

func (km KeyMap) Get(name string) *Binding {
	for _, bind := range km {
		for _, k := range bind.Keys() {
			if k == name {
				return bind
			}
		}
	}
	return km.New(name)
}

func (km KeyMap) New(keys ...string) *Binding {
	b := NewBinding(keys...)
	km.AddBind(b)
	return b
}

func (km KeyMap) AddBind(b *Binding) {
	km = append(km, b)
}

var Global = KeyMap{
	Quit(),
	ShowHelp(),
}

func DefaultListKeyMap() KeyMap {
	var km = KeyMap{
		Quit(),
		ToggleItem(),
		Up(),
		Down(),
		HalfPgUp(),
		HalfPgDown(),
		PgUp(),
		PgDown(),
		Home(),
		End(),
	}
	return km
}

func VimListKeyMap() KeyMap {
	var km = KeyMap{
		Quit(),
		ToggleItem().AddKeys(" "),
		Up().AddKeys("k"),
		Down().AddKeys("j"),
		HalfPgUp().AddKeys("K"),
		HalfPgDown().AddKeys("J"),
		PgUp(),
		PgDown(),
		Home().AddKeys("g"),
		End().AddKeys("G"),
	}
	return km
}

func Enter() *Binding {
	return NewBinding("enter")
}

func Filter() *Binding {
	return NewBinding("/").
		WithHelp("filter items")
}

func Save() *Binding {
	return NewBinding("ctrl+s").
		WithHelp("save edit")
}

func HalfPgUp() *Binding {
	return NewBinding("ctrl+u").
		WithHelp("½ page up").
		Cmd(HalfPageUp)
}

func HalfPgDown() *Binding {
	return NewBinding("ctrl+d").
		WithHelp("½ page down").
		Cmd(HalfPageDown)
}

func PgUp() *Binding {
	return NewBinding("pgup").
		WithHelp("page up").
		Cmd(PageUp)
}

func PgDown() *Binding {
	return NewBinding("pgdown").
		WithHelp("page down").
		Cmd(PageDown)
}

func End() *Binding {
	return NewBinding("end").
		WithHelp("list bottom").
		Cmd(Bottom)
}

func Home() *Binding {
	return NewBinding("home").
		WithHelp("list top").
		Cmd(Top)
}

func Up() *Binding {
	return NewBinding("up").
		WithHelp("move up").
		Cmd(LineUp)
}

func Down() *Binding {
	return NewBinding("down").
		WithHelp("move down").
		Cmd(LineDown)
}

func Next() *Binding {
	return NewBinding("right").
		WithHelp("next page").
		Cmd(NextPage)
}

func Prev() *Binding {
	return NewBinding("left").
		WithHelp("prev page").
		Cmd(PrevPage)
}

func ToggleItem() *Binding {
	return NewBinding("tab").
		WithHelp("select item").
		Cmd(message.ToggleItem())
}

func ToggleMatch() *Binding {
	return NewBinding("tab").
		WithHelp("select item")
}

func Quit() *Binding {
	return NewBinding("ctrl+c").
		WithHelp("quit program").
		Cmd(reactea.Destroy)
}

func ShowHelp() *Binding {
	return NewBinding("f1").
		WithHelp("help").
		Cmd(message.ShowHelp())
}

func Yes() *Binding {
	return NewBinding("y").
		WithHelp("confirm action")
}

func No() *Binding {
	return NewBinding("n").
		WithHelp("reject action")
}

func Esc() *Binding {
	return NewBinding("esc").
		WithHelp("exit screen")
}

func Edit() *Binding {
	return NewBinding("e").
		WithHelp("edit field")
}

type LineUpMsg struct{}
type HalfPageUpMsg struct{}
type PageUpMsg struct{}

func LineUp() tea.Msg     { return LineUpMsg{} }
func HalfPageUp() tea.Msg { return HalfPageUpMsg{} }
func PageUp() tea.Msg     { return PageUpMsg{} }

type LineDownMsg struct{}
type HalfPageDownMsg struct{}
type PageDownMsg struct{}

func LineDown() tea.Msg     { return LineDownMsg{} }
func HalfPageDown() tea.Msg { return HalfPageDownMsg{} }
func PageDown() tea.Msg     { return PageDownMsg{} }

type NextMsg struct{}
type PrevMsg struct{}

func NextPage() tea.Msg { return NextMsg{} }
func PrevPage() tea.Msg { return PrevMsg{} }

type TopMsg struct{}
type BottomMsg struct{}

func Top() tea.Msg    { return TopMsg{} }
func Bottom() tea.Msg { return BottomMsg{} }

type ToggleMsg struct {
	Index int
}

func Toggle(idx int) tea.Msg {
	return func() tea.Msg {
		return ToggleMsg{Index: idx}
	}
}
