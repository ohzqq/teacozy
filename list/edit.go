package list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/bubbles/key"
	"github.com/ohzqq/bubbles/list"
	"github.com/ohzqq/bubbles/textinput"
	"github.com/ohzqq/teacozy/util"
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

	filterInput := textinput.New()
	filterInput.Prompt = "New Item: "
	filterInput.PromptStyle = m.Styles.FilterPrompt
	filterInput.Cursor.Style = m.Styles.FilterCursor

	edit := &Edit{
		Model: m,
		input: filterInput,
	}

	var li []list.Item
	for _, i := range items() {
		li = append(li, i)
	}
	w, h := util.TermSize()

	del := edit.ItemDelegate()

	l := list.New(li, del, w, h)
	l.SetLimit(0)
	l.SetFilteringEnabled(false)

	edit.Model.Model = l

	return edit
}

func (e *Edit) ItemDelegate() list.DefaultDelegate {
	del := list.NewDefaultDelegate()
	del.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var cmds []tea.Cmd
		var cmd tea.Cmd

		switch msg := msg.(type) {
		case tea.KeyMsg:
			if e.input.Focused() {
				switch msg.Type {
				case tea.KeyEnter:
					item := NewItem(e.input.Value())
					cmd = m.InsertItem(m.Index()+1, item)
					cmds = append(cmds, cmd)
					e.input.Reset()
					e.input.Blur()
					m.SetHeight(m.Height() + 1)
				}
				e.input, cmd = e.input.Update(msg)
				cmds = append(cmds, cmd)
			} else {
				switch {
				case key.Matches(msg, keymap.Insert):
					m.SetHeight(m.Height() - 1)
					return e.input.Focus()
				case key.Matches(msg, keymap.Delete):
					m.RemoveItem(m.Index())
				}
			}
		}
		return nil
	}
	return del
}

func (m *Edit) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)

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
