package teacozy

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type Menus map[string]*Menu

func (m Menus) Get(key string) *Menu {
	return m[key]
}

func (m Menus) Set(key string, menu *Menu) Menus {
	m[key] = menu
	return m
}

type Menu struct {
	Model     *Info
	width     int
	Toggle    key.Binding
	height    int
	Label     string
	content   string
	show      bool
	style     lipgloss.Style
	IsFocused bool
	Items     []Key
}

func NewMenu(l string, toggle key.Binding, items ...Key) *Menu {
	m := DefaultMenu().SetKeys(items...)
	m.Label = l
	m.Toggle = toggle
	return m
}

func DefaultMenu() *Menu {
	m := Menu{
		Model: NewInfo(),
	}
	return &m
}

func (m Menu) Get(k string) Key {
	for _, item := range m.Items {
		if k == item.Key() {
			return item
		}
	}
	return Key{}
}

func (m Menu) Keys() []string {
	var keys []string
	for _, item := range m.Items {
		keys = append(keys, item.Key())
	}
	return keys
}

func (m *Menu) SetKeys(keys ...Key) *Menu {
	m.Items = keys
	for _, k := range keys {
		m.Model.Fields.Add(k)
	}
	return m
}

func (m *Menu) NewKey(k, h string, cmd MenuFunc) *Menu {
	key := NewMenuItem(k, h, cmd)
	m.AddKey(key)
	return m
}

func (m *Menu) AddKey(key Key) *Menu {
	m.Model.Fields.Add(key)
	m.Items = append(m.Items, key)
	return m
}

func (m *Menu) SetLabel(l string) *Menu {
	m.Label = l
	return m
}

func (m *Menu) SetToggle(toggle, help string) *Menu {
	m.Toggle = NewKeyBind(toggle, help)
	return m
}

func (m *Menu) View() string {
	return m.Model.View()
}

type Key struct {
	Bind key.Binding
	Cmd  MenuFunc
}

func NewMenuItem(k, h string, cmd MenuFunc) Key {
	return Key{
		Bind: NewKeyBind(k, h),
		Cmd:  cmd,
	}
}

func (i Key) Key() string {
	return i.Bind.Help().Key
}

func (i Key) Value() string {
	return i.Bind.Help().Desc
}

func (i Key) Set(v string) {}

func (i Key) String() string {
	return i.Bind.Help().Key + ": " + i.Bind.Help().Desc
}
