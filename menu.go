package teacozy

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Menus map[string]*Menu

func (m Menus) Get(key string) *Menu {
	return m[key]
}

type Menu struct {
	Model     viewport.Model
	width     int
	Toggle    key.Binding
	height    int
	Label     string
	content   string
	show      bool
	style     lipgloss.Style
	IsFocused bool
	Keys      []MenuItem
	Update    func(tea.Model, tea.Msg) tea.Cmd
}

type MenuItems []MenuItem

func NewMenu(l string, toggle key.Binding, items ...MenuItem) *Menu {
	m := Menu{
		Label:  l,
		Toggle: toggle,
		Keys:   items,
	}
	m.BuildModel()

	return &m
}

func UpdateMenu(m *TUI, msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.CurrentMenu.Toggle):
			m.HideMenu()
			cmds = append(cmds, SetFocusedViewCmd("list"))
		default:
			for _, item := range m.CurrentMenu.Keys {
				if key.Matches(msg, item.Key) {
					cmds = append(cmds, item.Cmd(m))
					m.HideMenu()
				}
			}
			m.HideMenu()
			cmds = append(cmds, SetFocusedViewCmd("list"))
		}
	}
	m.CurrentMenu.Model, cmd = m.CurrentMenu.Model.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m *Menu) SetKeys(keys MenuItems) *Menu {
	m.Keys = keys
	return m
}

func (m *Menu) BuildModel() {
	m.content = m.Render()
	vp := viewport.New(m.Width(), m.Height())
	vp.SetContent(m.content)
	m.Model = vp
}

func (m Menu) Render() string {
	var kh []string
	for _, k := range m.Keys {
		kh = append(kh, k.String())
	}
	style := FrameStyle().Copy().Width(m.Width())
	return style.Render(strings.Join(kh, "\n"))
}

func (m *Menu) SetWidth(w int) *Menu {
	m.width = w
	return m
}

func (m Menu) Width() int {
	if m.width != 0 {
		return m.width
	}
	return TermWidth() - 2
}

func (m Menu) Height() int {
	return lipgloss.Height(m.content)
}

type MenuItem struct {
	Key key.Binding
	Cmd MenuFunc
}

func NewMenuItem(k, h string, cmd MenuFunc) MenuItem {
	return MenuItem{
		Key: key.NewBinding(
			key.WithKeys(k),
			key.WithHelp(k, h),
		),
		Cmd: cmd,
	}
}

func (i MenuItem) String() string {
	return i.Key.Help().Key + ": " + i.Key.Help().Desc
}
