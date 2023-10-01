//go:build exclude

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

func NewEditableList(items Items) *Edit {
	m := New(items)

	edit := &Edit{
		Model: m,
	}

	del := list.NewDefaultDelegate()
	del.SetListType(list.Ol)

	edit.Model.SetDelegate(del)

	edit.input = edit.NewTextinputModel()
	edit.input.Prompt = "New Item: "

	edit.Model.SetLimit(1)
	edit.Model.SetFilteringEnabled(false)

	help := func() []key.Binding {
		return []key.Binding{
			edit.Model.KeyMap.InsertItem,
			edit.Model.KeyMap.RemoveItem,
		}
	}

	edit.Model.AdditionalShortHelpKeys = help
	edit.Model.AdditionalFullHelpKeys = help

	return edit
}

func (m *Edit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	if m.input.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEsc:
				m.resetInput()
			case tea.KeyEnter:
				item := NewItem(m.input.Value())
				cmd = m.InsertItem(m.Index()+1, item)
				cmds = append(cmds, cmd)
				m.resetInput()
			}
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		}
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.KeyMap.InsertItem):
				m.SetShowTitle(false)
				m.SetHeight(m.Height() - 1)
				m.state = Input
				return m, m.input.Focus()
			case key.Matches(msg, m.KeyMap.RemoveItem):
				m.RemoveItem(m.Index())
			}
		}
		m.Model, cmd = m.Model.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// ResetCommand resets the current filtering state.
func (m *Edit) ResetInput() {
	m.resetInput()
}

func (m *Edit) resetInput() {
	if m.state == Browsing {
		return
	}

	m.state = Browsing
	m.SetShowTitle(true)
	m.input.Reset()
	m.input.Blur()
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
