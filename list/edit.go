package list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/bubbles/key"
	"github.com/ohzqq/bubbles/list"
	"github.com/ohzqq/bubbles/textinput"
)

type Edit struct {
	*Model
	input textinput.Model
}

type KeyMap struct {
	Insert key.Binding
	Delete key.Binding
}

var keymap = KeyMap{
	Insert: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "insert"),
	),
	Delete: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "delete"),
	),
}

func NewEditableList(items Items) *Edit {
	m := New(items)

	edit := &Edit{
		Model: m,
	}

	edit.input = edit.NewTextinputModel()
	edit.input.Prompt = "New Item: "

	edit.Model.SetDelegate(edit.ItemDelegate())

	return edit
}

func (e *Edit) ItemDelegate() list.DefaultDelegate {
	del := list.NewDefaultDelegate()
	del.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keymap.Insert):
				m.SetHeight(m.Height() - 1)
				e.UpdateKeys(Input)
				return e.input.Focus()
			case key.Matches(msg, keymap.Delete):
				m.RemoveItem(m.Index())
			}
		}
		return nil
	}
	return del
}

func (m *Edit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if m.input.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				item := NewItem(m.input.Value())
				cmd = m.InsertItem(m.Index()+1, item)
				cmds = append(cmds, cmd)
				m.input.Reset()
				m.input.Blur()
				m.SetHeight(m.Height() + 1)
			}
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	} else {
		m.Model, cmd = m.Model.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (e Edit) View() string {
	var views []string

	if e.input.Focused() {
		in := e.input.View()
		views = append(views, in)
	}

	li := e.Model.View()
	views = append(views, li)

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}
