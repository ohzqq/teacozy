package teacozy

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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

func (m Menus) Del(key string) {
	delete(m, key)
}

type Menu struct {
	*Info
	Toggle    key.Binding
	Label     string
	content   string
	show      bool
	style     lipgloss.Style
	IsFocused bool
	Items     []*Key
}

func NewMenu(l string, toggle key.Binding, items ...*Key) *Menu {
	m := DefaultMenu().SetKeys(items...)
	m.Label = l
	m.Toggle = toggle
	return m
}

func DefaultMenu() *Menu {
	m := Menu{
		Info: NewInfo(),
	}
	return &m
}

func (m Menu) Get(k string) *Key {
	for _, item := range m.Items {
		if k == item.Key() {
			return item
		}
	}
	return &Key{}
}

func (m Menu) Keys() []string {
	var keys []string
	for _, item := range m.Items {
		keys = append(keys, item.Key())
	}
	return keys
}

func (m *Menu) SetKeys(keys ...*Key) *Menu {
	m.Items = keys
	for _, k := range keys {
		m.Info.Fields.Add(k)
	}
	return m
}

func (m *Menu) NewKey(k, h string, cmd MenuFunc) *Menu {
	key := NewKey(k, h).SetCmd(cmd)
	m.AddKey(key)
	return m
}

func (m *Menu) AddKey(key *Key) *Menu {
	m.Info.Fields.Add(key)
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

func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Toggle):
			m.show = false
			cmds = append(cmds, HideMenuCmd())
		default:
			for _, item := range m.Items {
				if key.Matches(msg, item.Binding) {
					m.show = false
					cmds = append(cmds, item.Cmd(m))
					cmds = append(cmds, HideMenuCmd())
				}
			}
			m.show = false
			cmds = append(cmds, HideMenuCmd())
		}
	}
	m.Info, cmd = m.Info.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Menu) Init() tea.Cmd { return nil }

func (m *Menu) View() string {
	return m.Info.View()
}
