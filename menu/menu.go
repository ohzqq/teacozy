package menu

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	urkey "github.com/ohzqq/teacozy/key"
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
	Keys      MenuItems
	Label     string
	Content   string
	isVisible bool
	height    int
	width     int
}

type MenuItems []Item

func NewMenu(l string, toggle key.Binding) *Menu {
	m := Menu{
		Label:  l,
		Toggle: toggle,
		Style:  style.FrameStyle(),
	}
	m.Model = viewport.New(m.Width(), m.Height())
	return &m
}

func (m *Menu) Init() tea.Cmd {
	m.Show()
	return UpdateMenuContentCmd(m.Render())
}
func (m *Menu) View() string {
	if m.isVisible {
		return m.Style.Render(m.Model.View())
	}
	return ""
}

func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, urkey.Quit):
			cmds = append(cmds, tea.Quit)
		case key.Matches(msg, m.Toggle):
			m.isVisible = !m.isVisible
			//m.Hide()
			//cmds = append(cmds, list.SetFocusedViewCmd("list"))
		default:
			for _, item := range m.Keys {
				if key.Matches(msg, item.Key) {
					cmds = append(cmds, item.Cmd())
					m.Hide()
				}
			}
			m.Hide()
			//cmds = append(cmds, list.SetFocusedViewCmd("list"))
		}
	case UpdateMenuContentMsg:
		m.Content = msg.Content
		m.Model.SetContent(m.Content)
	}
	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Menu) SetKeys(keys MenuItems) *Menu {
	m.Keys = keys
	return m
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
	//style := style.FrameStyle().Copy().Width(m.Width())
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
