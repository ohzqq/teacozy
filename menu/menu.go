package menu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
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
	*info.Info
	Toggle    *key.Key
	Label     string
	content   string
	show      bool
	style     lipgloss.Style
	IsFocused bool
	KeyMap    key.KeyMap
	Items     []*key.Key
}

func New(toggle, help string, keymap key.KeyMap) *Menu {
	m := Menu{
		KeyMap: keymap,
	}
	m.SetToggle(toggle, help)
	m.Info = info.New(keymap)
	return &m
}

func (m Menu) Get(k string) *key.Key {
	for _, item := range m.KeyMap.All() {
		if k == item.Key() {
			return item
		}
	}
	return &key.Key{}
}

func (m Menu) Keys() []string {
	var keys []string
	for _, item := range m.KeyMap.All() {
		keys = append(keys, item.Key())
	}
	return keys
}

func (m *Menu) AddKey(key *key.Key) *Menu {
	m.KeyMap.Add(key)
	return m
}

func (m *Menu) SetLabel(l string) *Menu {
	m.Label = l
	return m
}

func (m *Menu) SetToggle(toggle, help string) *Menu {
	//m.Toggle = key.NewBinding(
	//  key.WithKeys(toggle),
	//  key.WithHelp(toggle, help),
	//)
	m.Toggle = key.NewKey(toggle, help)
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
		case m.Toggle.Matches(msg):
			m.show = false
			cmds = append(cmds, HideMenuCmd())
		default:
			for _, item := range m.KeyMap.All() {
				if key.Matches(msg, item.Binding()) {
					m.show = false
					cmds = append(cmds, item.Cmd()(m))
					cmds = append(cmds, HideMenuCmd())
				}
			}
			m.show = false
			cmds = append(cmds, HideMenuCmd())
		}
	case UpdateMenuContentMsg:
		m.Info.SetContent(string(msg))
	}
	var i tea.Model
	i, cmd = m.Info.Update(msg)
	m.Info = i.(*info.Info)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Menu) Init() tea.Cmd { return nil }

func (m *Menu) View() string {
	return m.Info.View()
}
