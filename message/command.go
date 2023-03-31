package message

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

type StopEditingMsg struct{}
type StartEditingMsg struct{}
type SaveEditMsg struct{}
type StopFilteringMsg struct{}
type StartFilteringMsg struct{}
type ToggleItemMsg struct{}

func StartFilteringCmd() tea.Cmd {
	return func() tea.Msg {
		return StartFilteringMsg{}
	}
}

func StopFilteringCmd() tea.Cmd {
	return func() tea.Msg {
		return StopFilteringMsg{}
	}
}

func StopEditingCmd() tea.Cmd {
	return func() tea.Msg {
		return StopEditingMsg{}
	}
}

func SaveEditCmd() tea.Cmd {
	return func() tea.Msg {
		return SaveEditMsg{}
	}
}

func StartEditingCmd() tea.Cmd {
	return func() tea.Msg {
		return StartEditingMsg{}
	}
}

func ToggleItemCmd() tea.Cmd {
	return func() tea.Msg {
		return ToggleItemMsg{}
	}
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
