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

func (km KeyMap) Map() []map[string]string {
	c := make([]map[string]string, len(km))
	for i, k := range km {
		c[i] = map[string]string{k.Help().Key: k.Help().Desc}
	}
	return c
}

var Global = KeyMap{
	Quit(),
	ShowHelp(),
}

func Up() *Binding {
	return NewBinding("up").
		WithHelp("move up").
		Cmd(message.Up())
}

func Down() *Binding {
	return NewBinding("down").
		WithHelp("move down").
		Cmd(message.Down())
}

func Next() *Binding {
	return NewBinding("right").
		WithHelp("next page").
		Cmd(message.Next())
}

func Prev() *Binding {
	return NewBinding("left").
		WithHelp("prev page").
		Cmd(message.Prev())
}

func ToggleItem() *Binding {
	return NewBinding("tab").
		WithHelp("select item").
		Cmd(message.ToggleItem())
}

func Quit() *Binding {
	return NewBinding("ctrl+c").
		WithHelp("quit program").
		Cmd(message.Quit())
}

func ShowHelp() *Binding {
	return NewBinding("H").
		WithHelp("help").
		Cmd(message.ShowHelp())
}
