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
	Frame  style.Frame
	funcs  map[string]MenuFunc
	Label  string
	Toggle *key.Key
}

func NewMenu(k, h string) *Menu {
	m := Menu{
		Info:   info.New(),
		Label:  h,
		Toggle: key.NewKey(k, h),
		KeyMap: key.NewKeyMap(),
		funcs:  make(map[string]MenuFunc),
	}
	m.Frame = DefaultWidgetStyle()
	return &m
}

func (m *Menu) Add(k, h string, cmd MenuFunc) *Menu {
	m.New(k, h)
	m.funcs[h] = cmd
	return m
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
					cmds = append(cmds, HideMenuCmd())
					cmds = append(cmds, fn(ui))
				}
			}
			cmds = append(cmds, HideMenuCmd())
			cmds = append(cmds, SetFocusedViewCmd("list"))
		}
	}
	var model tea.Model
	model, cmd = m.Info.Update(msg)
	cmds = append(cmds, cmd)

	m.Info = model.(*info.Info)

	return tea.Batch(cmds...)
}

func (m *Menu) SetKeyMap(km key.KeyMap) *Menu {
	m.KeyMap = km
	return m
}

func (m *Menu) SetToggle(toggle, help string) *Menu {
	m.Label = help
	m.Toggle = key.NewKey(toggle, help)
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
