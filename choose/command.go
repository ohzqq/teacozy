package choose

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/style"
	"golang.org/x/exp/maps"
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

func (m *Choose) Header(text string) *Choose {
	m.header = text
	return m
}

func (m *Choose) ChoiceMap(choices []map[string]string) *Choose {
	m.choiceMap = choices
	return m
}

func (m *Choose) SetStyle(s style.List) *Choose {
	m.Style = s
	return m
}

func (m *Choose) Limit(l int) *Choose {
	m.limit = l
	return m
}

func (m *Choose) NoLimit() *Choose {
	return m.Limit(len(m.Choices))
}

func (m *Choose) SetHeight(h int) *Choose {
	m.Height = h
	return m
}

func (m *Choose) SetWidth(w int) *Choose {
	m.Width = w
	return m
}

func (m *Choose) SetSize(w, h int) *Choose {
	m.SetWidth(w)
	m.SetHeight(h)
	return m
}

type ReturnSelectionsMsg struct{}

func ReturnSelectionsCmd(m *Choose) tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
	}
}

func QuitCmd(m *Choose) tea.Cmd {
	return func() tea.Msg {
		m.quitting = true
		return ReturnSelectionsMsg{}
	}
}

func FilterItemsCmd(m *Choose) tea.Cmd {
	return func() tea.Msg {
		return nil
	}
}

type StopFilteringMsg struct{}
type StartFilteringMsg struct{}

func StartFilteringCmd(m *Choose) tea.Cmd {
	return func() tea.Msg {
		return StartFilteringMsg{}
	}
}

func SelectItemCmd(m *Choose) tea.Cmd {
	return func() tea.Msg {
		if m.limit == 1 {
			return nil
		}
		m.ToggleSelection()
		return nil
	}
}

func UpCmd(m *Choose) tea.Cmd {
	return func() tea.Msg {
		m.CursorUp()
		return nil
	}
}

func DownCmd(m *Choose) tea.Cmd {
	return func() tea.Msg {
		m.CursorDown()
		return nil
	}
}

func TopCmd(m *Choose) tea.Cmd {
	return func() tea.Msg {
		m.Cursor = 0
		m.Paginator.Page = 0
		return nil
	}
}

func BottomCmd(m *Choose) tea.Cmd {
	return func() tea.Msg {
		m.Cursor = len(m.Items.Items) - 1
		m.Paginator.Page = m.Paginator.TotalPages - 1
		return nil
	}
}

func NextPageCmd(m *Choose) tea.Cmd {
	return func() tea.Msg {
		m.Cursor = clamp(0, len(m.Items.Items)-1, m.Cursor+m.Height)
		m.Paginator.NextPage()
		return nil
	}
}

func PrevPageCmd(m *Choose) tea.Cmd {
	return func() tea.Msg {
		m.Cursor = clamp(0, len(m.Items.Items)-1, m.Cursor-m.Height)
		m.Paginator.PrevPage()
		return nil
	}
}

func SelectAllItemsCmd(m *Choose) tea.Cmd {
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

func DeselectAllItemsCmd(m *Choose) tea.Cmd {
	return func() tea.Msg {
		if m.limit <= 1 {
			return nil
		}

		maps.Clear(m.Selected)
		m.numSelected = 0

		return nil
	}
}

func FReturnSelectionsCmd(m *Filter) tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
	}
}

func FQuitCmd(m *Filter) tea.Cmd {
	return func() tea.Msg {
		m.quitting = true
		return ReturnSelectionsMsg{}
	}
}

func FStopFilteringCmd(m *Filter) tea.Cmd {
	return func() tea.Msg {
		if m.limit == 1 {
			m.ToggleSelection()
			return ReturnSelectionsMsg{}
		}

		m.Input.Reset()
		m.Input.Blur()
		return StopFilteringMsg{}
	}
}

func FSelectItemCmd(m *Filter) tea.Cmd {
	return func() tea.Msg {
		if m.limit == 1 {
			return nil
		}
		m.ToggleSelection()
		return nil
	}
}

func FUpCmd(m *Filter) tea.Cmd {
	return func() tea.Msg {
		m.CursorUp()
		return nil
	}
}

func FDownCmd(m *Filter) tea.Cmd {
	return func() tea.Msg {
		m.CursorDown()
		return nil
	}
}
