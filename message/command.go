package message

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ReturnSelectionsMsg struct{}

func ReturnSelections() tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
	}
}

type QuitMsg struct{}

func Quit() tea.Cmd {
	return func() tea.Msg {
		return QuitMsg{}
	}
}

type ChangeRouteMsg struct {
	Name string
}

type GetConfirmationMsg struct {
	Question string
}

func GetConfirmation(q string) tea.Cmd {
	return func() tea.Msg {
		return GetConfirmationMsg{Question: q}
	}
}

type ConfirmMsg struct {
	Confirmed bool
}

func Confirm(confirm bool) tea.Cmd {
	return func() tea.Msg {
		return ConfirmMsg{Confirmed: confirm}
	}
}

type StopEditingMsg struct{}
type StartEditingMsg struct{}
type StopFilteringMsg struct{}
type StartFilteringMsg struct{}
type ToggleItemMsg struct{}
type SaveEditMsg struct{}

func StartFiltering() tea.Cmd {
	return func() tea.Msg {
		return StartFilteringMsg{}
	}
}

func StopFiltering() tea.Cmd {
	return func() tea.Msg {
		return StopFilteringMsg{}
	}
}

func ChangeRoute(name string) tea.Cmd {
	return func() tea.Msg {
		return ChangeRouteMsg{Name: name}
	}
}

func StopEditing() tea.Cmd {
	return func() tea.Msg {
		return StopEditingMsg{}
	}
}

func SaveEdit() tea.Cmd {
	return func() tea.Msg {
		return SaveEditMsg{}
	}
}

func StartEditing() tea.Cmd {
	return func() tea.Msg {
		return StartEditingMsg{}
	}
}

func ToggleItem() tea.Cmd {
	return func() tea.Msg {
		return ToggleItemMsg{}
	}
}

type ShowHelpMsg struct{}

func ShowHelp() tea.Cmd {
	return func() tea.Msg {
		return ShowHelpMsg{}
	}
}

type HideHelpMsg struct{}

func HideHelp() tea.Cmd {
	return func() tea.Msg {
		return HideHelpMsg{}
	}
}

type UpMsg struct{}
type DownMsg struct{}
type NextMsg struct{}
type PrevMsg struct{}
type TopMsg struct{}
type BottomMsg struct{}

func Up() tea.Cmd {
	return func() tea.Msg {
		return UpMsg{}
	}
}

func Down() tea.Cmd {
	return func() tea.Msg {
		return DownMsg{}
	}
}

func Next() tea.Cmd {
	return func() tea.Msg {
		return NextMsg{}
	}
}

func Prev() tea.Cmd {
	return func() tea.Msg {
		return PrevMsg{}
	}
}

func Top() tea.Cmd {
	return func() tea.Msg {
		return TopMsg{}
	}
}

func Bottom() tea.Cmd {
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
