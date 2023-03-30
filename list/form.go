package list

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/londek/reactea"
	"github.com/ohzqq/teacozy/item"
	"github.com/ohzqq/teacozy/style"
)

type Form struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[FormProps]
	Cursor      int
	Matches     []item.Item
	Input       textarea.Model
	Viewport    *viewport.Model
	quitting    bool
	Placeholder string
	Prompt      string
	Style       style.List
}

type FormProps struct {
	Props
	EditItem func(int)
}

type FormKeys struct {
	Up          key.Binding
	Down        key.Binding
	ToggleItem  key.Binding
	Quit        key.Binding
	StopEditing key.Binding
	SaveValue   key.Binding
	Edit        key.Binding
}

func NewField() *Field {
	f := &Field{}
	return f
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
		m.Cursor = clamp(0, len(m.Props().Visible())-1, m.Cursor+1)
		if m.Cursor >= m.Viewport.YOffset+m.Viewport.Height {
			m.Viewport.LineDown(1)
		}
	case StartEditingMsg:
		reactea.SetCurrentRoute("field")
		return nil
	case StopEditingMsg:
		m.Input.Reset()
		m.Input.Blur()
		reactea.SetCurrentRoute("default")
		return nil
	case tea.KeyMsg:
		if m.Input.Focused() {
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
			case key.Matches(msg, formKey.Edit):
				if m.Input.Focused() {
					//m.Props().SetValue(cur.idx, m.Input.Value())
					m.Input.Blur()
				} else {
					//m.Props().EditItem(m.Cursor)
					m.Input.SetValue(m.Props().Visible()[m.Cursor].Str)
					return m.Input.Focus()
					//cmds = append(cmds, StartEditingCmd())
					//cmds = append(cmds, EditItemCmd(cur))
				}
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

	v := viewport.New(0, 0)
	tm.Viewport = &v
	//tm.Input.Focus()

	return nil
}

type FieldKeys struct {
	Exit key.Binding
	Save key.Binding
	Quit key.Binding
}

type Field struct {
	reactea.BasicComponent
	reactea.BasicPropfulComponent[FieldProps]

	Input textarea.Model
	Val   string
}

type FieldProps struct {
	item.Item
	Save func(int, string)
}

func (m *Field) Init(props FieldProps) tea.Cmd {
	m.UpdateProps(props)
	m.Input = textarea.New()
	m.Input.ShowLineNumbers = false
	m.Input.SetValue(props.Str)
	return m.Input.Focus()
}

func (m *Field) Render(w, h int) string {
	m.Input.SetWidth(w)

	return m.Input.View()
}

func (m *Field) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case StopEditingMsg:
		m.Input.Blur()
		reactea.SetCurrentRoute("form")
		return nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, fieldKey.Quit):
			return reactea.Destroy
		case key.Matches(msg, fieldKey.Exit):
			cmds = append(cmds, StopEditingCmd())
		case key.Matches(msg, fieldKey.Save):
			i := m.Props().Item
			i.Write(m.Input.Value())
			cmds = append(cmds, StopEditingCmd())
		}
		m.Input, cmd = m.Input.Update(msg)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}
