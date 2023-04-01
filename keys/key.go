package keys

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/message"
)

type KeyMap []*Binding

type Binding struct {
	key.Binding
	help   string
	keys   []string
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
	k.keys = keys
	k.WithHelp(k.help)
	return k
}

func (k *Binding) WithHelp(h string) *Binding {
	k.help = h
	k.Binding.SetHelp(strings.Join(k.keys, "/"), h)
	return k
}

var Global = KeyMap{
	Quit(),
	ShowHelp(),
}

func Up() *Binding {
	return NewBinding("up").
		WithHelp("move up").
		Cmd(message.UpCmd())
}

func Down() *Binding {
	return NewBinding("down").
		WithHelp("move down").
		Cmd(message.DownCmd())
}

func Next() *Binding {
	return NewBinding("right").
		WithHelp("next page").
		Cmd(message.NextCmd())
}

func Prev() *Binding {
	return NewBinding("left").
		WithHelp("prev page").
		Cmd(message.PrevCmd())
}

func ToggleItem() *Binding {
	return NewBinding("tab").
		WithHelp("select item").
		Cmd(message.ToggleItemCmd())
}

func Quit() *Binding {
	return NewBinding("ctrl+c").
		WithHelp("quit program").
		Cmd(message.QuitCmd())
}

func ShowHelp() *Binding {
	return NewBinding("ctrl+h").
		WithHelp("help").
		Cmd(message.ShowHelpCmd())
}
