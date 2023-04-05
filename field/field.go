package field

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

type Field struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[Props]

	Input    textarea.Model
	quitting bool
	Prompt   string
}

type KeyMap struct {
	Quit        key.Binding
	StopEditing key.Binding
	Save        key.Binding
	Edit        key.Binding
}

func (m Field) KeyMap() keys.KeyMap {
	km := keys.KeyMap{
		keys.NewBinding("esc").
			WithHelp("stop editing").
			Cmd(StopEditing()),
		keys.NewBinding("e").
			WithHelp("edit field").
			Cmd(StartEditing()),
		keys.NewBinding("ctrl+s").
			WithHelp("save edits").
			Cmd(m.SaveEdit()),
	}
	return km
}

type Props struct {
	*props.Item
	fields string
}

func New() *Field {
	tm := Field{
		Prompt: style.PromptPrefix,
	}
	return &tm
}

func NewProps(i *props.Item, fields string) Props {
	return Props{
		Item:   i,
		fields: fields,
	}
}

func (c Field) Initializer(props *props.Items) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := New()
		p := NewProps(props.CurrentItem(), props.Snapshot)
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
	case StopEditingMsg:
		m.Input.Blur()
		cmds = append(cmds, message.ChangeRoute("default"))
	case message.ConfirmMsg:
		if msg.Confirmed {
			cmds = append(cmds, SaveEdit())
		}
		cmds = append(cmds, StopEditing())
	case SaveEditMsg:
		m.Props().Item.Str = m.Input.Value()
		m.Input.Reset()
		cmds = append(cmds, StopEditing())
	case StartEditingMsg:
		textarea.Blink()
		return m.Input.Focus()
	case tea.KeyMsg:
		if m.Input.Focused() {
			for _, k := range m.KeyMap() {
				if key.Matches(msg, k.Binding) {
					m.Input.Blur()
					cmds = append(cmds, k.TeaCmd)
				}
			}
			m.Input, cmd = m.Input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			for _, k := range m.KeyMap() {
				if key.Matches(msg, k.Binding) {
					cmds = append(cmds, k.TeaCmd)
				}
			}
		}
	}

	return tea.Batch(cmds...)
}

func (m *Field) SaveEdit() tea.Cmd {
	if m.Props().Item.Str != m.Input.Value() {
		return message.GetConfirmation("Save edit?")
	}
	return StopEditing()
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

type StopEditingMsg struct{}

func StopEditing() tea.Cmd {
	return func() tea.Msg {
		return StopEditingMsg{}
	}
}

type SaveEditMsg struct{}

func SaveEdit() tea.Cmd {
	return func() tea.Msg {
		return SaveEditMsg{}
	}
}

type StartEditingMsg struct{}

func StartEditing() tea.Cmd {
	return func() tea.Msg {
		return StartEditingMsg{}
	}
}
