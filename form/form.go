package form

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/londek/reactea/router"
	"github.com/ohzqq/teacozy/message"
	"github.com/ohzqq/teacozy/props"
	"github.com/ohzqq/teacozy/style"
)

type Form struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[FormProps]
	Cursor      int
	Input       textarea.Model
	Viewport    *viewport.Model
	quitting    bool
	Placeholder string
	Prompt      string
	Style       style.List
}

type FormProps struct {
	*props.Items
	Save func([]map[string]string)
}

type FormKeys struct {
	Up          key.Binding
	Down        key.Binding
	ToggleItem  key.Binding
	Quit        key.Binding
	StopEditing key.Binding
	Save        key.Binding
	Edit        key.Binding
}

func NewForm() *Form {
	tm := Form{
		Style:  style.ListDefaults(),
		Prompt: style.PromptPrefix,
	}
	return &tm
}

func FormRouteInitializer(props FormProps) router.RouteInitializer {
	return func(router.Params) (reactea.SomeComponent, tea.Cmd) {
		component := NewForm()
		return component, component.Init(props)
	}
}

func (m *Form) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case message.ReturnSelectionsMsg:
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	case message.UpMsg:
		m.Cursor = clamp(0, len(m.Props().Visible())-1, m.Cursor-1)
		if m.Cursor < m.Viewport.YOffset {
			m.Viewport.SetYOffset(m.Cursor)
		}
	case message.DownMsg:
		h := lipgloss.Height(m.Props().Visible()[m.Cursor].Str)
		m.Cursor = clamp(0, len(m.Props().Visible())-1, m.Cursor+1)
		if m.Cursor >= m.Viewport.YOffset+m.Viewport.Height-h {
			m.Viewport.LineDown(h)
		} else if m.Cursor == len(m.Props().Visible())-1 {
			m.Viewport.GotoBottom()
		}
	case message.StopEditingMsg:
		m.Input.Reset()
		m.Input.Blur()
		reactea.SetCurrentRoute("default")
		return nil
	case tea.KeyMsg:
		if m.Input.Focused() {
			switch {
			case key.Matches(msg, formKey.Save):
				m.Props().Visible()[m.Cursor].Str = m.Input.Value()
				m.Input.Blur()
			}
			m.Input, cmd = m.Input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case key.Matches(msg, formKey.StopEditing):
				cmds = append(cmds, message.StopEditingCmd())
			case key.Matches(msg, formKey.Up):
				cmds = append(cmds, message.UpCmd())
			case key.Matches(msg, formKey.Down):
				cmds = append(cmds, message.DownCmd())
			case key.Matches(msg, formKey.ToggleItem):
				cmds = append(cmds, message.ToggleItemCmd())
			case key.Matches(msg, formKey.Quit):
				m.quitting = true
				cmds = append(cmds, message.ReturnSelectionsCmd())
			case key.Matches(msg, formKey.Save):
				m.Props().Save(m.Props().Items.Map())
				reactea.SetCurrentRoute("default")
				return nil
			case key.Matches(msg, formKey.Edit):
				m.Input.SetValue(m.Props().Visible()[m.Cursor].Str)
				return m.Input.Focus()
			}
		}
	}

	m.Cursor = clamp(0, len(m.Props().Visible())-1, m.Cursor)
	return tea.Batch(cmds...)
}

func (m *Form) Render(w, h int) string {
	m.Viewport.Height = h
	m.Viewport.Width = w

	var s strings.Builder
	items := m.Props().RenderItems(m.Cursor, m.Props().Visible())
	s.WriteString(items)

	m.Viewport.SetContent(s.String())

	view := m.Viewport.View()
	if m.Input.Focused() {
		view += "\n" + m.Input.View()
	}

	return view
}

func (tm *Form) Init(props FormProps) tea.Cmd {
	tm.UpdateProps(props)

	tm.Input = textarea.New()
	tm.Input.Prompt = tm.Prompt
	tm.Input.Placeholder = tm.Placeholder
	tm.Input.SetWidth(tm.Props().Width)
	tm.Input.ShowLineNumbers = false

	v := viewport.New(0, 0)
	tm.Viewport = &v
	//tm.Input.Focus()

	return nil
}

func clamp(min, max, val int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

var formKey = FormKeys{
	Save: key.NewBinding(
		key.WithKeys("ctrl+w"),
		key.WithHelp("ctrl+w", "save"),
	),
	ToggleItem: key.NewBinding(
		key.WithKeys(" ", "tab"),
		key.WithHelp("space", "select item"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("down", "move cursor down"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("up", "move cursor up"),
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
