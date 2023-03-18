package keymap

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap []Binding
type Binding struct {
	key.Binding
	opts []key.BindingOpt
	Cmd  tea.Cmd
}

type BindingOpt func(*Binding)

func Matches(k tea.KeyMsg, b ...Binding) bool {
	var binds []key.Binding
	for _, kb := range b {
		binds = append(binds, kb.Binding)
	}
	return key.Matches(k, binds...)
}

func NewBinding(opts ...BindingOpt) Binding {
	k := Binding{}
	for _, opt := range opts {
		opt(&k)
	}
	k.Binding = key.NewBinding(k.opts...)
	return k
}

func (k *Binding) SetCmd(cmd tea.Cmd) {
	k.Cmd = cmd
}

func (k Binding) Matches(msg tea.KeyMsg) bool {
	return key.Matches(msg, k.Binding)
}

func WithKeys(keys ...string) BindingOpt {
	return func(b *Binding) {
		b.opts = append(b.opts, key.WithKeys(keys...))
	}
}

func WithHelp(k, desc string) BindingOpt {
	return func(b *Binding) {
		b.opts = append(b.opts, key.WithHelp(k, desc))
	}
}

func WithCmd(cmd tea.Cmd) BindingOpt {
	return func(b *Binding) {
		b.Cmd = cmd
	}
}
