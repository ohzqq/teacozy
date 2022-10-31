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

type Field struct {
	Model   textarea.Model
	width   int
	height  int
	label   string
	toggle  key.Binding
	content string
	focus   bool
	style   lipgloss.Style
}

func (f Field) Toggle() key.Binding { return f.toggle }
func (f Field) Label() string       { return f.label }
func (f Field) Focused() bool       { return f.Model.Focused() }

func (f *Field) Focus() tea.Cmd {
	f.Model.Focus()
	return nil
}

func (f *Field) Blur() {
	f.Model.Blur()
}

func NewField(title, content string) *Field {
	field := Field{
		label:   title,
		content: content,
		height:  lipgloss.Height(content),
		width:   util.TermWidth() - 4,
	}

	t := key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit field"),
	)
	field.toggle = t
	field.Model = textarea.New()
	field.SetValue(content)
	field.SetHeight(field.height)
	field.SetWidth(field.width)
	field.ShowLineNumbers = false

	return &field
}

func (m Field) Update(list *list.List, msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, cozykey.SaveAndExit) {
			m.SetContent(m.Model.Value())
			m.Blur()
		}
		if m.Focused() {
			m.Model, cmd = m.Model.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
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
