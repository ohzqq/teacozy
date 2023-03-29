package filter

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/style"
)

func DefaultStyle() style.List {
	var s style.List
	s.Cursor = style.Cursor
	s.SelectedPrefix = style.Selected
	s.UnselectedPrefix = style.Unselected
	s.Text = style.Foreground
	s.Match = lipgloss.NewStyle().Foreground(color.Cyan())
	s.Header = lipgloss.NewStyle().Foreground(color.Purple())
	s.Prompt = style.Prompt
	return s
}

func (m *Filter) Header(text string) *Filter {
	m.header = text
	return m
}

func (m *Filter) ChoiceMap(choices []map[string]string) *Filter {
	m.choiceMap = choices
	return m
}

func (m *Filter) SetStyle(s style.List) *Filter {
	m.Style = s
	return m
}

func (m *Filter) Limit(l int) *Filter {
	m.limit = l
	return m
}

func (m *Filter) NoLimit() *Filter {
	return m.Limit(len(m.Choices))
}

func (m *Filter) SetHeight(h int) *Filter {
	m.Height = h
	return m
}

func (m *Filter) SetWidth(w int) *Filter {
	m.Width = w
	return m
}

func (m *Filter) SetSize(w, h int) *Filter {
	m.SetWidth(w)
	m.SetHeight(h)
	return m
}

type ReturnSelectionsMsg struct{}

func ReturnSelectionsCmd(m *Filter) tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
	}
}

func QuitCmd(m *Filter) tea.Cmd {
	return func() tea.Msg {
		m.quitting = true
		return ReturnSelectionsMsg{}
	}
}

func StopFilteringCmd(m *Filter) tea.Cmd {
	return func() tea.Msg {
		if m.limit == 1 {
			m.ToggleSelection()
			return ReturnSelectionsMsg{}
		}

		m.Input.Reset()
		m.Input.Blur()
		return nil
	}
}

func SelectItemCmd(m *Filter) tea.Cmd {
	return func() tea.Msg {
		if m.limit == 1 {
			return nil
		}
		m.ToggleSelection()
		return nil
	}
}

func UpCmd(m *Filter) tea.Cmd {
	return func() tea.Msg {
		m.CursorUp()
		return nil
	}
}

func DownCmd(m *Filter) tea.Cmd {
	return func() tea.Msg {
		m.CursorDown()
		return nil
	}
}
