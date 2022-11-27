package menu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
)

type Menus map[string]*Menu

type Menu struct {
	*info.Info
	key.KeyMap
	Label  string
	Toggle *key.Key
}

func New(toggle, help string) *Menu {
	m := Menu{
		Info:   info.New(),
		KeyMap: key.NewKeyMap(),
		Label:  help,
	}
	m.SetToggle(toggle, help)
	return &m
}

func (m *Menu) SetKeyMap(km key.KeyMap) *Menu {
	m.KeyMap = km
	return m
}

//func (m *Menu) AddKey(key *key.Key) *Menu {
//  m.KeyMap.Add(key)
//  return m
//}

func (m *Menu) SetToggle(toggle, help string) *Menu {
	m.Label = help
	m.Toggle = key.NewKey(toggle, help)
	return m
}

func (m Menu) Get(k string) *key.Key {
	for _, item := range m.KeyMap.All() {
		if k == item.Name() {
			return item
		}
	}
	return &key.Key{}
}

func (m Menu) Keys() []string {
	var keys []string
	for _, item := range m.KeyMap.All() {
		keys = append(keys, item.Name())
	}
	return keys
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
			m.Info.Toggle()
		default:
			for _, item := range m.KeyMap.All() {
				if key.Matches(msg, item.Binding()) {
					m.Hide()
					cmds = append(cmds, item.Cmd()(m))
					cmds = append(cmds, HideMenuCmd())
				}
			}
			m.Hide()
			cmds = append(cmds, HideMenuCmd())
		}
	case UpdateMenuContentMsg:
		m.Info.SetContent(string(msg))
	}

	var i tea.Model
	i, cmd = m.Info.Update(msg)
	//m.Info, cmd = m.Info.Update(msg)
	cmds = append(cmds, cmd)

	m.Info = i.(*info.Info)

	return m, tea.Batch(cmds...)
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m *Menu) View() string {
	m.Info.NewSection().SetTitle(m.Label).SetFields(m.KeyMap)
	return m.Info.View()
}

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
