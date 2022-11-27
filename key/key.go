package key

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Key struct {
	key key.Binding
}

func Matches(msg tea.KeyMsg, bind ...key.Binding) bool {
	return key.Matches(msg, bind...)
}

func NewKey(k, h string) *Key {
	return &Key{
		key: key.NewBinding(
			key.WithKeys(k),
			key.WithHelp(k, h),
		),
	}
}

func NewBinding(k, help string) key.Binding {
	return key.NewBinding(
		key.WithKeys(k),
		key.WithHelp(k, help),
	)
}

func (k Key) Binding() key.Binding {
	return k.key
}

func (k Key) Matches(msg tea.KeyMsg) bool {
	return key.Matches(msg, k.key)
}

func (i Key) Name() string {
	return i.key.Help().Key
}

func (i Key) Content() string {
	return i.key.Help().Desc
}

func (i Key) Set(v string) {}

func (i Key) String() string {
	return i.key.Help().Key + ": " + i.key.Help().Desc
}
