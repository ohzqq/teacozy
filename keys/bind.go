package keys

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

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

func (k *Binding) Disable() {
	k.Binding.SetEnabled(false)
}

func (k *Binding) Enable() {
	k.Binding.SetEnabled(true)
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

func Enter() *Binding {
	return New("enter")
}

func Filter() *Binding {
	return New("/").
		WithHelp("filter items")
	//Cmd(ChangeRoute("filter"))
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
	//Cmd(UpdateItem(ToggleItems))
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
		Cmd(ChangeRoute("prev"))
}

func Edit() *Binding {
	return New("e").
		WithHelp("edit field").
		Cmd(EditItem)
}

func Save() *Binding {
	return New("ctrl+s").
		WithHelp("save edit")
}
