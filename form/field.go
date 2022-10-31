package form

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/tea/key"
	cozykey "github.com/ohzqq/teacozy/key"
	"github.com/ohzqq/teacozy/list"
	"github.com/ohzqq/teacozy/util"
)

type Fields map[string]*Field

func (m Fields) Get(key string) *Field {
	return m[key]
}

type Field struct {
	Model     textarea.Model
	width     int
	height    int
	title     string
	content   string
	show      bool
	style     lipgloss.Style
	IsFocused bool
	Update    func(tea.Model, tea.Msg) tea.Cmd
}

func NewField(title, content string) *Field {
	field := Field{
		title:   title,
		content: content,
		height:  lipgloss.Height(content),
		width:   util.TermWidth() - 4,
	}

	field.Model = textarea.New()
	field.SetValue(content)
	field.SetHeight(field.height)
	field.SetWidth(field.width)
	field.ShowLineNumbers = false

	return &field
}

func (m Field) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, cozykey.SaveAndExit) {
			m.SetContent(m.Model.Value())
			m.Model.Blur()
		}
		if m.Model.Focused() {
			m.Model, cmd = m.Model.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

type UpdateFieldContentMsg string

func (f *Field) UpdateFieldContentCmd(value string) tea.Cmd {
	return func() tea.Msg {
		f.content = value
		return UpdateFieldContentMsg(content)
	}
}

func (f *Field) SetContent(content string) {
	f.content = content
	f.Model.SetValue(content)
}

func UpdateField(m *list.List, msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, cozykey.SaveAndExit) {
			m.SetContent(m.Model.Value())
			m.Model.Blur()
		}
		if m.Model.Focused() {
			m.Model, cmd = m.Model.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	m.CurrentField.Model, cmd = m.CurrentField.Model.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m Field) View() string {
	return m.Model.View()
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
