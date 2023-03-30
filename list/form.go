package list

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
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
	Props
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
		Style:  DefaultStyle(),
		Prompt: style.PromptPrefix,
	}
	return &tm
}

func (m *Form) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.Props().Height == 0 || m.Props().Height > msg.Height {
			m.Viewport.Height = msg.Height - lipgloss.Height(m.Input.View())
		}

		m.Viewport.Width = msg.Width
	case ReturnSelectionsMsg:
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	case UpMsg:
		m.Cursor = clamp(0, len(m.Props().Visible())-1, m.Cursor-1)
		if m.Cursor < m.Viewport.YOffset {
			m.Viewport.SetYOffset(m.Cursor)
		}
	case DownMsg:
		h := lipgloss.Height(m.Props().Visible()[m.Cursor].Str)
		m.Cursor = clamp(0, len(m.Props().Visible())-1, m.Cursor+1)
		if m.Cursor >= m.Viewport.YOffset+m.Viewport.Height-h {
			m.Viewport.LineDown(h)
		} else if m.Cursor == len(m.Props().Visible())-1 {
			m.Viewport.GotoBottom()
		}
	case StopEditingMsg:
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
				cmds = append(cmds, StopEditingCmd())
			case key.Matches(msg, formKey.Up):
				cmds = append(cmds, UpCmd())
			case key.Matches(msg, formKey.Down):
				cmds = append(cmds, DownCmd())
			case key.Matches(msg, formKey.ToggleItem):
				cmds = append(cmds, ToggleItemCmd())
			case key.Matches(msg, formKey.Quit):
				m.quitting = true
				cmds = append(cmds, ReturnSelectionsCmd())
			case key.Matches(msg, formKey.Save):
				m.Props().Save(m.Props().Items.Map())
				reactea.SetCurrentRoute("default")
				return nil
			case key.Matches(msg, formKey.Edit):
				m.Input.SetValue(m.Props().Visible()[m.Cursor].Str)
				return m.Input.Focus()
				//cmds = append(cmds, StartEditingCmd())
				//cmds = append(cmds, EditItemCmd(cur))
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
