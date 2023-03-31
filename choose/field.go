package choose

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

type Field struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[FieldProps]
	Input       textarea.Model
	quitting    bool
	Placeholder string
	Prompt      string
}

type FieldKeyMap struct {
	Quit        key.Binding
	StopEditing key.Binding
	Save        key.Binding
	Edit        key.Binding
}

type FieldProps struct {
	*props.Item
	Save func(map[string]string)
}

func NewField() *Field {
	tm := Field{
		Prompt: style.PromptPrefix,
	}
	return &tm
}

func (m *Field) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case message.StopEditingMsg:
		m.Input.Reset()
		m.Input.Blur()
		reactea.SetCurrentRoute("default")
		return nil
	case message.SaveEditMsg:
		m.Input.Blur()
	case message.StartEditingMsg:
		textarea.Blink()
		return m.Input.Focus()
	case tea.KeyMsg:
		if m.Input.Focused() {
			switch {
			case key.Matches(msg, formKey.Save):
				m.Input.Blur()
			}
			m.Input, cmd = m.Input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case key.Matches(msg, formKey.StopEditing):
				cmds = append(cmds, message.StopEditingCmd())
			case key.Matches(msg, formKey.Quit):
				m.quitting = true
			case key.Matches(msg, formKey.Save):
				m.Props().Item.Str = m.Input.Value()
				cmds = append(cmds, message.StopEditingCmd())
			case key.Matches(msg, formKey.Edit):
				cmds = append(cmds, message.StartEditingCmd())
			}

		}
	}

	return tea.Batch(cmds...)
}

func (m *Field) Render(w, h int) string {
	m.Input.SetWidth(w)
	return m.Input.View()
}

func (tm *Field) Init(props FieldProps) tea.Cmd {
	tm.UpdateProps(props)

	tm.Input = textarea.New()
	tm.Input.Prompt = tm.Prompt
	tm.Input.Placeholder = tm.Placeholder
	tm.Input.ShowLineNumbers = false

	tm.Input.SetValue(tm.Props().Item.Str)

	textarea.Blink()
	return tm.Input.Focus()
}

var formKey = FieldKeyMap{
	Save: key.NewBinding(
		key.WithKeys("ctrl+w"),
		key.WithHelp("ctrl+w", "save"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	StopEditing: key.NewBinding(
		key.WithKeys("esc", "q"),
		key.WithHelp("esc/q", "stop editing"),
	),
	Edit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "edit field"),
	),
}
