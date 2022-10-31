package list

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/util"
)

type Fields map[string]*Field

func (m Fields) Get(key string) *Field {
	return m[key]
}

type Field struct {
	Model     textarea.Model
	width     int
	Toggle    key.Binding
	height    int
	Label     string
	content   string
	show      bool
	style     lipgloss.Style
	IsFocused bool
	Keys      FieldItems
	Update    func(tea.Model, tea.Msg) tea.Cmd
}

type FieldItems []FieldItem

func NewField(l string, toggle key.Binding) *Field {
	return &Field{
		Label:  l,
		Toggle: toggle,
	}
}

func UpdateField(m *List, msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.CurrentMenu.Toggle):
			m.ShowTextArea = false
			cmds = append(cmds, SetFocusedViewCmd("list"))
		default:
			for _, item := range m.CurrentMenu.Keys {
				if key.Matches(msg, item.Key) {
					cmds = append(cmds, item.Cmd(m))
					m.ShowTextArea = false
				}
			}
			m.ShowTextArea = false
			cmds = append(cmds, SetFocusedViewCmd("list"))
		}
	}
	m.CurrentMenu.Model, cmd = m.CurrentMenu.Model.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m *Field) SetKeys(keys FieldItems) *Field {
	m.Keys = keys
	return m
}

func (m *Field) BuildModel() {
	m.content = m.Render()
	area := textarea.New()
	area.SetValue(m.content)
	m.Model = area
}

func (m Field) Render() string {
	var kh []string
	for _, k := range m.Keys {
		kh = append(kh, k.String())
	}
	style := FrameStyle().Copy().Width(m.Width())
	return style.Render(strings.Join(kh, "\n"))
}

func (m *Field) SetWidth(w int) *Field {
	m.width = w
	return m
}

func (m Field) Width() int {
	if m.width != 0 {
		return m.width
	}
	return util.TermWidth() - 2
}

func (m Field) Height() int {
	return lipgloss.Height(m.content)
}

type FieldItem struct {
	Key key.Binding
	Cmd MenuCmd
}

func NewFieldItem(k, h string, cmd MenuCmd) FieldItem {
	return FieldItem{
		Key: key.NewBinding(
			key.WithKeys(k),
			key.WithHelp(k, h),
		),
		Cmd: cmd,
	}
}

func (i FieldItem) String() string {
	return i.Key.Help().Key + ": " + i.Key.Help().Desc
}
