package list

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/util"
)

type Menu struct {
	Model   viewport.Model
	width   int
	toggle  key.Binding
	height  int
	label   string
	content string
	show    bool
	focus   bool
	style   lipgloss.Style
	Keys    []MenuItem
}

func (m Menu) Label() string {
	return m.label
}

func NewMenu(l string, toggle key.Binding, keys []MenuItem) *Menu {
	m := Menu{
		label:  l,
		toggle: toggle,
		Keys:   keys,
	}

	m.content = m.Render()
	vp := viewport.New(m.Width(), m.Height())
	vp.SetContent(m.content)
	m.Model = vp

	return &m
}

func (m Menu) Update(list *List, msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Toggle()):
			cmds = append(cmds, FocusListCmd())
		default:
			for _, item := range m.Keys {
				if key.Matches(msg, item.Key) {
					cmds = append(cmds, item.Cmd(list))
					list.HideWidget()
				}
			}
			cmds = append(cmds, FocusListCmd())
		}
	}
	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m Menu) Toggle() key.Binding {
	return m.toggle
}

func (m *Menu) Blur() {
	m.focus = false
}

func (m *Menu) SetKeys(keys []MenuItem) *Menu {
	m.Keys = keys
	return m
}

func (m Menu) Render() string {
	var kh []string
	for _, k := range m.Keys {
		kh = append(kh, k.String())
	}
	return strings.Join(kh, "\n")
}

func (m Menu) View() string {
	style := FrameStyle().Copy().Width(m.Width())
	return style.Render(m.Render())
}

func (m *Menu) SetWidth(w int) *Menu {
	m.width = w
	return m
}

func (m *Menu) SetContent(content string) {
	m.content = content
	m.Model.SetContent(content)
}

func (m Menu) Width() int {
	if m.width != 0 {
		return m.width
	}
	return util.TermWidth() - 2
}

func (m Menu) Height() int {
	return lipgloss.Height(m.content)
}

func (m *Menu) Focus() tea.Cmd {
	m.focus = true
	return nil
}

func (m Menu) Focused() bool {
	return m.focus
}

type MenuCmd func(m *List) tea.Cmd

type MenuItem struct {
	Key key.Binding
	Cmd MenuCmd
}

func NewMenuItem(k, h string, cmd MenuCmd) MenuItem {
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
