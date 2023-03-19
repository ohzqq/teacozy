package list

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/style"
	"golang.org/x/exp/maps"
)

func (m *Model) Header(text string) *Model {
	m.header = text
	return m
}

//func (m *Model) Style(s ListStyle) *Model {
//  m.ListStyle = s
//  return m
//}

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

func DefaultStyle() Style {
	var s Style
	s.Cursor = style.Cursor
	s.SelectedPrefix = style.Selected
	s.UnselectedPrefix = style.Unselected
	s.Text = style.Foreground
	s.Match = lipgloss.NewStyle().Foreground(color.Cyan())
	s.Header = lipgloss.NewStyle().Foreground(color.Purple())
	s.Prompt = style.Prompt
	return s
}

func FilterItemsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.filterState = Filtering
		m.textinput.Focus()
		return textinput.Blink()
	}
}

func StopFilteringCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.filterState = Unfiltered
		m.textinput.Reset()
		m.textinput.Blur()
		return nil
	}
}

func (m Model) Chosen() []string {
	var chosen []string
	if len(m.selected) > 0 {
		for k := range m.selected {
			chosen = append(chosen, m.Choices[k])
		}
	} else if len(m.matches) > m.cursor && m.cursor >= 0 {
		chosen = append(chosen, m.matches[m.cursor].Str)
	}
	return chosen
}

type ReturnSelectionsMsg struct{}

func ReturnSelectionsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
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
		m.paginator.Page = 0
		return nil
	}
}

func BottomCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.cursor = len(m.items) - 1
		m.paginator.Page = m.paginator.TotalPages - 1
		return nil
	}
}

func NextPageCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.cursor = clamp(0, len(m.items)-1, m.cursor+m.height)
		m.paginator.NextPage()
		return nil
	}
}

func PrevPageCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.cursor = clamp(0, len(m.items)-1, m.cursor-m.height)
		m.paginator.PrevPage()
		return nil
	}
}

func SelectAllItemsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.limit <= 1 {
			return nil
		}
		for i := range m.matches {
			if m.numSelected >= m.limit {
				break // do not exceed given limit
			}
			if _, ok := m.selected[i]; ok {
				continue
			} else {
				m.selected[m.matches[i].Index] = struct{}{}
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

		maps.Clear(m.selected)
		m.numSelected = 0

		return nil
	}
}
