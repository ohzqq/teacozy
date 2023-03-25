package list

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/style"
	"golang.org/x/exp/maps"
)

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

func (m *Model) Height(h int) *Model {
	m.height = h
	return m
}

func (m *Model) Width(w int) *Model {
	m.width = w
	return m
}

func (m *Model) SetSize(w, h int) *Model {
	m.Width(w)
	m.Height(h)
	return m
}

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

func FilterItemsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.filterState = Filtering
		m.Input.Focus()
		return textinput.Blink()
	}
}

func StopFilteringCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.limit == 1 {
			m.ToggleSelection()
			return ReturnSelectionsMsg{}
		}

		m.filterState = Unfiltered
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

func TopCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.cursor = 0
		m.Paginator.Page = 0
		return nil
	}
}

func BottomCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.cursor = len(m.Items) - 1
		m.Paginator.Page = m.Paginator.TotalPages - 1
		return nil
	}
}

func NextPageCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.cursor = clamp(0, len(m.Items)-1, m.cursor+m.height)
		m.Paginator.NextPage()
		return nil
	}
}

func PrevPageCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.cursor = clamp(0, len(m.Items)-1, m.cursor-m.height)
		m.Paginator.PrevPage()
		return nil
	}
}

func SelectAllItemsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.limit <= 1 {
			return nil
		}
		for i := range m.Matches {
			if m.numSelected >= m.limit {
				break // do not exceed given limit
			}
			if _, ok := m.Selected[i]; ok {
				continue
			} else {
				m.Selected[m.Matches[i].Index] = struct{}{}
				m.numSelected++
			}
		}
		return nil
	}
}

func DeselectAllItemsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.limit <= 1 {
			return nil
		}

		maps.Clear(m.Selected)
		m.numSelected = 0

		return nil
	}
}
