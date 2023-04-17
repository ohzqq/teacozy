package keys

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
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
	Help(),
}

func Enter() *Binding {
	return NewBinding("enter")
}

func Filter() *Binding {
	return NewBinding("/").
		WithHelp("filter items").
		Cmd(ChangeRoute("filter"))
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
		Cmd(func() tea.Msg { return ToggleItemMsg{} })
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

func Help() *Binding {
	return NewBinding("f1").
		WithHelp("show help")
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
		WithHelp("exit screen").
		Cmd(ReturnToList)
}

func Edit() *Binding {
	return NewBinding("e").
		WithHelp("edit field").
		Cmd(ChangeRoute("edit"))
}
