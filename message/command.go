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

type ChangeRouteMsg struct {
	Name string
}

type StopEditingMsg struct{}
type StartEditingMsg struct{}
type StopFilteringMsg struct{}
type StartFilteringMsg struct{}
type ToggleItemMsg struct{}
type SaveEditMsg struct{}

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

func ChangeRouteCmd(name string) tea.Cmd {
	return func() tea.Msg {
		return ChangeRouteMsg{Name: name}
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
type DownMsg struct{}
type NextMsg struct{}
type PrevMsg struct{}
type TopMsg struct{}
type BottomMsg struct{}

func UpCmd() tea.Cmd {
	return func() tea.Msg {
		return UpMsg{}
	}
}

func DownCmd() tea.Cmd {
	return func() tea.Msg {
		return DownMsg{}
	}
}

func NextCmd() tea.Cmd {
	return func() tea.Msg {
		return NextMsg{}
	}
}

func PrevCmd() tea.Cmd {
	return func() tea.Msg {
		return PrevMsg{}
	}
}

func TopCmd() tea.Cmd {
	return func() tea.Msg {
		return TopMsg{}
	}
}

func BottomCmd() tea.Cmd {
	return func() tea.Msg {
		return BottomMsg{}
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
