package keys

import (
	"github.com/charmbracelet/bubbles/key"
	"golang.org/x/exp/slices"
)

type KeyMap struct {
	keys []*Binding
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

func (km *KeyMap) Replace(b *Binding) {
	if idx := km.Index(b); idx != -1 {
		km.keys = slices.Replace(km.Keys(), idx, idx+1, b)
	}
}

func (km KeyMap) String(idx int) string {
	return km.keys[idx].Help().Desc
}

func (km KeyMap) Label(idx int) string {
	return km.keys[idx].Help().Key
}

func (km KeyMap) Len() int {
	return len(km.keys)
}

func (km KeyMap) Set(int, string) {}

func MapKeys(keys ...key.Binding) []map[string]string {
	c := make([]map[string]string, len(keys))
	for i, k := range keys {
		c[i] = map[string]string{k.Help().Key: k.Help().Desc}
	}
	return c
}

func DefaultKeyMap() KeyMap {
	k := []*Binding{
		PgUp(),
		PgDown(),
		Up(),
		Down(),
		HalfPgUp(),
		HalfPgDown(),
		Home(),
		End(),
	}
	return NewKeyMap(k...)
}

func VimKeyMap() KeyMap {
	km := []*Binding{
		Up().AddKeys("k"),
		Down().AddKeys("j"),
		HalfPgUp().AddKeys("K"),
		HalfPgDown().AddKeys("J"),
		Home().AddKeys("g"),
		End().AddKeys("G"),
		Quit().AddKeys("q"),
	}
	return NewKeyMap(km...)
}
