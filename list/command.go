package list

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/maps"
)

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
