package keys

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/slices"
)

type KeyMap struct {
	keys []*Binding
}

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

func NewKeyMap(b ...*Binding) KeyMap {
	return KeyMap{
		keys: b,
	}
}

func (km KeyMap) Map() []map[string]string {
	c := make([]map[string]string, len(km.Keys()))
	for i, k := range km.Keys() {
		c[i] = map[string]string{k.Help().Key: k.Help().Desc}
	}
	return c
}

func (km *KeyMap) Get(name string) *Binding {
	for _, bind := range km.Keys() {
		for _, k := range bind.Keys() {
			if k == name {
				return bind
			}
		}
	}
	return km.New(name)
}

func (km KeyMap) Keys() []*Binding {
	return km.keys
}

func (km *KeyMap) New(keys ...string) *Binding {
	b := New(keys...)
	km.AddBinds(b)
	return b
}

func (km *KeyMap) AddBinds(b ...*Binding) {
	km.keys = append(km.keys, b...)
}

func (km KeyMap) Contains(bind *Binding) bool {
	return slices.ContainsFunc(km.Keys(), func(b *Binding) bool {
		for _, k := range bind.Keys() {
			return slices.Contains(b.Keys(), k)
		}
		return false
	})
}

func (km KeyMap) Index(bind *Binding) int {
	return slices.IndexFunc(km.Keys(), func(b *Binding) bool {
		return km.Contains(bind)
	})
}

func MapKeys(keys ...key.Binding) []map[string]string {
	c := make([]map[string]string, len(keys))
	for i, k := range keys {
		c[i] = map[string]string{k.Help().Key: k.Help().Desc}
	}
	return c
}
