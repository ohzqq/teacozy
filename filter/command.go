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

func (m Model) Chosen() []int {
	var chosen []int
	if m.quitting {
		return chosen
	} else if len(m.Selected) > 0 {
		for k := range m.Selected {
			chosen = append(chosen, k)
		}
	} else if len(m.Matches) > m.cursor && m.cursor >= 0 {
		chosen = append(chosen, m.cursor)
	}
	return chosen
}

func (m *Model) Header(text string) *Model {
	m.header = text
	return m
}

func (m *Model) ChoiceMap(choices []map[string]string) *Model {
	m.choiceMap = choices
	return m
}

func (m *Model) SetStyle(s style.List) *Model {
	m.Style = s
	return m
}

func (m *Model) Limit(l int) *Model {
	m.limit = l
	return m
}

func (m *Model) NoLimit() *Model {
	return m.Limit(len(m.Choices))
}

func (m *Model) SetHeight(h int) *Model {
	m.Height = h
	return m
}

func (m *Model) SetWidth(w int) *Model {
	m.Width = w
	return m
}

func (m *Model) SetSize(w, h int) *Model {
	m.SetWidth(w)
	m.SetHeight(h)
	return m
}

type ReturnSelectionsMsg struct{}

func ReturnSelectionsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
	}
}

func QuitCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.quitting = true
		return ReturnSelectionsMsg{}
	}
}

func StopFilteringCmd(m *Model) tea.Cmd {
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

func SelectItemCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.limit == 1 {
			return nil
		}
		m.ToggleSelection()
		return nil
	}
}

func UpCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.CursorUp()
		return nil
	}
}

func DownCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.CursorDown()
		return nil
	}
}
