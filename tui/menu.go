package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/style"
)

type Menus map[string]*Menu

type MenuFunc func(tea.Model) tea.Cmd

type Menu struct {
	*info.Info
	key.KeyMap
	Frame style.Frame
	funcs map[string]MenuFunc
	Label string
}

func MainMenu() *Menu {
	toggle := key.NewKey("m", "menu")
	m := NewMenu(toggle)
	m.SetSize(m.Frame.Width(), m.Frame.Height())
	m.NewKey("?", "help", GoToHelp)
	m.NewSection().SetTitle("Main menu").SetFields(m.KeyMap)
	return m
}

func NewMenu(toggle *key.Key) *Menu {
	m := Menu{
		Info:   info.New(),
		Label:  toggle.Content(),
		KeyMap: key.NewKeyMap(),
		funcs:  make(map[string]MenuFunc),
	}
	//for _, name := range m.Keys() {
	//  kb := m.GetKey(name)
	//}
	m.Info.Toggle = toggle
	m.Frame = DefaultWidgetStyle()
	return &m
}

func (m *Menu) NewKey(k, h string, cmd MenuFunc) *Menu {
	m.KeyMap.NewKey(k, h)
	m.funcs[h] = cmd
	return m
}

func (m *Menu) AddKey(k *key.Key, cmd MenuFunc) *Menu {
	m.KeyMap.Add(k)
	m.funcs[k.Content()] = cmd
	return m
}

func (m *Tui) updateMenu(msg tea.Msg) tea.Cmd {
	var (
		//cmd  tea.Cmd
		cmds []tea.Cmd
		//model tea.Model
	)

	switch msg := msg.(type) {
	//case info.HideInfoMsg:
	//case info.UpdateContentMsg:
	case tea.KeyMsg:
		for _, name := range m.CurrentMenu.Keys() {
			if kb := m.CurrentMenu.GetKey(name); kb.Matches(msg) {
				fn := m.CurrentMenu.funcs[name]
				cmds = append(cmds, fn(m))
				//cmds = append(cmds, HideMenuCmd())
				cmds = append(cmds, info.HideInfoCmd())
			}
			cmds = append(cmds, info.HideInfoCmd())
			//cmds = append(cmds, HideMenuCmd())
		}
	}
	//cmd = m.updateInfo(msg, m.CurrentMenu.Info)
	//cmds = append(cmds, cmd)

	cmds = append(cmds, info.UpdateContentCmd(m.CurrentMenu.Render()))
	return tea.Batch(cmds...)
}

func (m Menu) GetInfo() *info.Info {
	m.NewSection().SetTitle("opts").SetFields(m.KeyMap)
	return m.Info
}

func (m *Menu) Update(ui *TUI, msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case m.Toggle.Matches(msg):
			cmds = append(cmds, HideMenuCmd())
			cmds = append(cmds, SetFocusedViewCmd("list"))
		default:
			for _, name := range m.Keys() {
				if kb := m.GetKey(name); kb.Matches(msg) {
					fn := m.funcs[name]
					cmds = append(cmds, fn(ui))
					cmds = append(cmds, HideMenuCmd())
				}
			}
			cmds = append(cmds, HideMenuCmd())
			cmds = append(cmds, SetFocusedViewCmd("list"))
		}
	}

	var model tea.Model
	model, cmd = m.Info.Update(msg)
	m.Info = model.(*info.Info)
	//m.Info, cmd = m.Info.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Menu) View() string {
	return m.Info.View()
}

func (m *Menu) SetKeyMap(km key.KeyMap) *Menu {
	m.KeyMap = km
	return m
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

type UpdateMenuContentMsg string

func UpdateMenuContentCmd(s string) tea.Cmd {
	return func() tea.Msg {
		return UpdateMenuContentMsg(s)
	}
}

type HideMenuMsg struct{}

func HideMenuCmd() tea.Cmd {
	return func() tea.Msg {
		return HideMenuMsg{}
	}
}

type ShowMenuMsg struct{ *Menu }

func ShowMenuCmd(menu *Menu) tea.Cmd {
	return func() tea.Msg {
		return ShowMenuMsg{Menu: menu}
	}
}

type ChangeMenuMsg struct{ *Menu }

func GoToMenuCmd(menu *Menu) MenuFunc {
	return func(m tea.Model) tea.Cmd {
		return ChangeMenuCmd(menu)
	}
}

func ChangeMenuCmd(menu *Menu) tea.Cmd {
	return func() tea.Msg {
		return ChangeMenuMsg{Menu: menu}
	}
}
