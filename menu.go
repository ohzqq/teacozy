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
	Items     []MenuItem
}

func NewMenu(l string, toggle key.Binding, items ...MenuItem) *Menu {
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

func (m Menu) Get(k string) MenuItem {
	for _, item := range m.Items {
		if k == item.Key() {
			return item
		}
	}
	return MenuItem{}
}

func (m Menu) Keys() []string {
	var keys []string
	for _, item := range m.Items {
		keys = append(keys, item.Key())
	}
	return keys
}

func (m *Menu) SetKeys(keys ...MenuItem) *Menu {
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

func (m *Menu) AddKey(key MenuItem) *Menu {
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

type MenuItem struct {
	KeyBind key.Binding
	Cmd     MenuFunc
}

func NewMenuItem(k, h string, cmd MenuFunc) MenuItem {
	return MenuItem{
		KeyBind: NewKeyBind(k, h),
		Cmd:     cmd,
	}
}

func (i MenuItem) Key() string {
	return i.KeyBind.Help().Key
}

func (i MenuItem) Value() string {
	return i.KeyBind.Help().Desc
}

func (i MenuItem) Set(v string) {}

func (i MenuItem) String() string {
	return i.KeyBind.Help().Key + ": " + i.KeyBind.Help().Desc
}
