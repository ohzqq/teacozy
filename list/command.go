package list

import (
	tea "github.com/charmbracelet/bubbletea"
)

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
