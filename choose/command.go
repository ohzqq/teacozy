package choose

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

func ReturnSelectionsCmd() tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
	}
}

func QuitCmd() tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
	}
}

func FilterItemsCmd() tea.Cmd {
	return func() tea.Msg {
		return nil
	}
}

type StopFilteringMsg struct{}
type StartFilteringMsg struct{}
type ToggleItemMsg struct{}

func StartFilteringCmd(m *Choose) tea.Cmd {
	return func() tea.Msg {
		return StartFilteringMsg{}
	}
}

func StopFilteringCmd() tea.Cmd {
	return func() tea.Msg {
		return StopFilteringMsg{}
	}
}

func ToggleItemCmd() tea.Cmd {
	return func() tea.Msg {
		return ToggleItemMsg{}
	}
}

type SetCursorMsg struct {
	cursor int
}

type UpMsg struct{}

func UpCmd() tea.Cmd {
	return func() tea.Msg {
		return UpMsg{}
	}
}

type DownMsg struct{}

func DownCmd() tea.Cmd {
	return func() tea.Msg {
		return DownMsg{}
	}
}

//func SelectAllItemsCmd(m *Choose) tea.Cmd {
//  return func() tea.Msg {
//    if m.limit <= 1 {
//      return nil
//    }
//    for i := range m.Matches {
//      if m.numSelected >= m.limit {
//        break // do not exceed given limit
//      }
//      if _, ok := m.Selected[i]; ok {
//        continue
//      } else {
//        m.Selected[m.Matches[i].Index] = struct{}{}
//        m.numSelected++
//      }
//    }
//    return nil
//  }
//}

//func DeselectAllItemsCmd(m *Choose) tea.Cmd {
//  return func() tea.Msg {
//    if m.limit <= 1 {
//      return nil
//    }

//    maps.Clear(m.Selected)
//    m.numSelected = 0

//    return nil
//  }
//}
