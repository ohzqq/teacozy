package field

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

type Field struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Input       textarea.Model
	quitting    bool
	Placeholder string
	Prompt      string
}

type KeyMap struct {
	Quit        key.Binding
	StopEditing key.Binding
	Save        key.Binding
	Edit        key.Binding
}

type Props struct {
	*props.Item
	fields string
}

func NewField() *Field {
	tm := Field{
		Prompt: style.PromptPrefix,
	}
	return &tm
}

func NewFieldProps(i *props.Item, fields string) Props {
	return Props{
		Item:   i,
		fields: fields,
	}
}

func (c Field) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := NewField()
		p := NewFieldProps(props.CurrentItem(), props.Snapshot)
		return component, component.Init(p)
	}
}

func (c Field) Name() string {
	return "editField"
}

func (m *Field) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case message.StopEditingMsg:
		m.Input.Blur()
		cmds = append(cmds, message.ChangeRoute("choose"))
	case message.ConfirmMsg:
		if msg.Confirmed {
			cmds = append(cmds, message.SaveEdit())
		}
		cmds = append(cmds, message.StopEditing())
	case message.SaveEditMsg:
		m.Props().Item.Str = m.Input.Value()
		m.Input.Reset()
		cmds = append(cmds, message.StopEditing())
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
				cmds = append(cmds, message.StopEditing())
			case key.Matches(msg, formKey.Quit):
				m.quitting = true
			case key.Matches(msg, formKey.Save):
				if m.Props().Item.Str != m.Input.Value() {
					return message.GetConfirmation("Save edit?")
				}
				cmds = append(cmds, message.StopEditing())
			case key.Matches(msg, formKey.Edit):
				cmds = append(cmds, message.EditField())
			}
		}
	}

	return tea.Batch(cmds...)
}

func (m *Field) Render(w, h int) string {
	m.Input.SetWidth(w)
	lh := m.Props().Item.LineHeight()
	m.Input.SetHeight(lh)
	return lipgloss.JoinVertical(lipgloss.Left, m.Props().fields, m.Input.View())
}

func (tm *Field) Init(props Props) tea.Cmd {
	tm.UpdateProps(props)

	tm.Input = textarea.New()
	tm.Input.Prompt = tm.Prompt
	tm.Input.Placeholder = tm.Placeholder
	tm.Input.ShowLineNumbers = false

	tm.Input.SetValue(tm.Props().Item.Str)

	textarea.Blink()
	return tm.Input.Focus()
}

var formKey = KeyMap{
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
