package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/teacozy/info"
	"github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list"
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

func NewMenu(toggle *key.Key) *Menu {
	m := Menu{
		Info:   info.New(),
		Label:  toggle.Content(),
		KeyMap: key.NewKeyMap(),
		funcs:  make(map[string]MenuFunc),
	}
	m.Info.Toggle = toggle
	m.Frame = DefaultWidgetStyle()
	return &m
}

func MainMenu() *Menu {
	toggle := key.NewKey("m", "menu")
	m := NewMenu(toggle)
	m.SetSize(m.Frame.Width(), m.Frame.Height())
	m.AddKey(KeyMap().GetKey("?"), GoToHelpView)
	m.NewSection().SetTitle("Main menu").SetFields(m.KeyMap)
	return m
}

func ActionMenu() *Menu {
	toggle := key.NewKey("a", "action")
	m := NewMenu(toggle)
	m.SetSize(m.Frame.Width(), m.Frame.Height())
	m.NewKey("P", "print", PrintItemsMenuFunc)
	m.NewSection().SetTitle("Action menu").SetFields(m.KeyMap)
	return m
}

func EditItemMetaMenu(item *list.Item) *Menu {
	toggle := key.NewKey("i", "action")
	m := NewMenu(toggle)
	//item.Meta.Info.Title = "Meta (e: edit)"
	m.Info = item.Meta.Info
	m.SetSize(m.Frame.Width(), m.Frame.Height())
	m.NewKey("e", "edit", EditItemMeta(item))
	return m
}

func EditItemMeta(item *list.Item) MenuFunc {
	fn := func(m tea.Model) tea.Cmd {
		return EditItemMetaCmd(item)
	}
	return fn
}

func PrintItems(items ...*list.Item) tea.Cmd {
	for _, i := range items {
		fmt.Println(i.Content())
	}
	return tea.Quit
}

func PrintItemsMenuFunc(m tea.Model) tea.Cmd {
	ui := m.(*Tui)
	main := ui.Main.(*list.List)
	ui.state = mainModel
	main.SetAction(PrintItems)
	return list.ReturnSelectionsCmd()
}

func (m *Menu) NewKey(k, h string, cmd MenuFunc) *Menu {
	m.KeyMap.NewKey(k, h)
	m.funcs[k] = cmd
	return m
}

func (m *Menu) AddKey(k *key.Key, cmd MenuFunc) *Menu {
	m.KeyMap.Add(k)
	m.funcs[k.Name()] = cmd
	return m
}

func (m *Tui) updateMenu(msg tea.Msg) tea.Cmd {
	var (
		//cmd  tea.Cmd
		cmds []tea.Cmd
		//model tea.Model
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, name := range m.CurrentMenu.Keys() {
			if kb := m.CurrentMenu.GetKey(name); kb.Matches(msg) {
				fn := m.CurrentMenu.funcs[name]
				return fn(m)
			}
			cmds = append(cmds, info.HideInfoCmd())
		}
	}

	cmds = append(cmds, info.UpdateContentCmd(m.CurrentMenu.Render()))
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
