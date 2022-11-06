package menu

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
)

type Menus map[string]*Menu

func (m Menus) Get(key string) *Menu {
	return m[key]
}

type Menu struct {
	Model     viewport.Model
	Toggle    key.Binding
	Style     lipgloss.Style
	Keys      []Item
	Label     string
	Content   string
	isVisible bool
	height    int
	width     int
}

func NewMenu(l string, toggle key.Binding, keys ...Item) *Menu {
	m := Menu{
		Label:  l,
		Toggle: toggle,
		Style:  style.FrameStyle(),
		Keys:   keys,
	}
	m.Content = m.Render()
	m.Model = viewport.New(m.Width(), m.Height())
	return &m
}

func (m *Menu) Init() tea.Cmd {
	return UpdateMenuContentCmd(m.Render())
}
func (m *Menu) View() string {
	m.Model.SetContent(m.Content)
	if m.isVisible {
		return m.Style.Render(m.Model.View())
	}
	return ""
}

//func (m *Menu) Update(msg tea.Msg) (*Menu, tea.Cmd) {
func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			cmds = append(cmds, tea.Quit)
		}
		switch {
		case key.Matches(msg, m.Toggle):
			m.ToggleVisibility()
		default:
			for _, item := range m.Keys {
				if key.Matches(msg, item.Key) {
					cmds = append(cmds, item.Cmd())
					m.Hide()
				}
			}
			m.Hide()
		}
	case UpdateMenuContentMsg:
		m.Content = msg.Content
	}
	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Menu) SetKeys(keys ...Item) *Menu {
	m.Keys = keys
	return m
}

func (m *Menu) ToggleVisibility() {
	m.isVisible = !m.isVisible
}

func (m *Menu) Hide() {
	m.isVisible = false
}

func (m *Menu) Show() {
	m.isVisible = true
}

func (m Menu) Render() string {
	var kh []string
	for _, k := range m.Keys {
		kh = append(kh, k.String())
	}
	return strings.Join(kh, "\n")
}

func (m *Menu) SetWidth(w int) *Menu {
	m.width = w
	return m
}

func (m Menu) Width() int {
	if m.width != 0 {
		return m.width
	}
	return util.TermWidth() - 2
}

func (m Menu) Height() int {
	return lipgloss.Height(m.Content)
}

type Item struct {
	Key key.Binding
	Cmd CmdFunc
}

func NewItem(k, h string, cmd CmdFunc) Item {
	return Item{
		Key: key.NewBinding(
			key.WithKeys(k),
			key.WithHelp(k, h),
		),
		Cmd: cmd,
	}
}

func (i *Item) SetCmd(cmd CmdFunc) *Item {
	i.Cmd = cmd
	return i
}

func (i Item) String() string {
	return i.Key.Help().Key + ": " + i.Key.Help().Desc
}
