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
	TeaCmd tea.Cmd
}

func New(keys ...string) *Binding {
	k := Binding{
		Binding: key.NewBinding(),
	}
	k.WithKeys(keys...)
	return &k
}

func NewBind(k key.Binding) *Binding {
	return &Binding{
		Binding: k,
	}
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
	b := New(keys...)
	km.AddBind(b)
	return b
}

func (km KeyMap) AddBind(b *Binding) {
	km = append(km, b)
}

func MapKeys(keys ...key.Binding) []map[string]string {
	c := make([]map[string]string, len(keys))
	for i, k := range keys {
		c[i] = map[string]string{k.Help().Key: k.Help().Desc}
	}
	return c
}

var Global = KeyMap{
	Quit(),
}

func Enter() *Binding {
	return New("enter")
}

func Filter() *Binding {
	return New("/").
		WithHelp("filter items").
		Cmd(ChangeRoute("filter"))
}

func Save() *Binding {
	return New("ctrl+s").
		WithHelp("save edit")
}

func HalfPgUp() *Binding {
	return New("ctrl+u").
		WithHelp("½ page up").
		Cmd(HalfPageUp)
}

func HalfPgDown() *Binding {
	return New("ctrl+d").
		WithHelp("½ page down").
		Cmd(HalfPageDown)
}

func PgUp() *Binding {
	return New("pgup").
		WithHelp("page up").
		Cmd(PageUp)
}

func PgDown() *Binding {
	return New("pgdown").
		WithHelp("page down").
		Cmd(PageDown)
}

func End() *Binding {
	return New("end").
		WithHelp("list bottom").
		Cmd(Bottom)
}

func Home() *Binding {
	return New("home").
		WithHelp("list top").
		Cmd(Top)
}

func Up() *Binding {
	return New("up").
		WithHelp("move up").
		Cmd(LineUp)
}

func Down() *Binding {
	return New("down").
		WithHelp("move down").
		Cmd(LineDown)
}

func Next() *Binding {
	return New("right").
		WithHelp("next page").
		Cmd(NextPage)
}

func Prev() *Binding {
	return New("left").
		WithHelp("prev page").
		Cmd(PrevPage)
}

func Toggle() *Binding {
	return New("tab").
		WithHelp("select item").
		Cmd(ToggleItem)
}

func Quit() *Binding {
	return New("ctrl+c").
		WithHelp("quit program").
		Cmd(reactea.Destroy)
}

func Help() *Binding {
	return New("f1").
		WithHelp("show help").
		Cmd(ShowHelp)
}

func Yes() *Binding {
	return New("y").
		WithHelp("confirm action")
}

func No() *Binding {
	return New("n").
		WithHelp("reject action")
}

func Esc() *Binding {
	return New("esc").
		WithHelp("exit screen").
		Cmd(ReturnToList)
}

func Edit() *Binding {
	return New("e").
		WithHelp("edit field").
		Cmd(ChangeRoute("edit"))
}
