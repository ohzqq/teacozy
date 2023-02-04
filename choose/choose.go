// Package choose provides an interface to choose one option from a given list
// of options. The options can be provided as (new-line separated) stdin or a
// list of arguments.
//
// It is different from the filter command as it does not provide a fuzzy
// finding input, so it is best used for smaller lists of options.
//
// Let's pick from a list of gum flavors:
//
// $ gum choose "Strawberry" "Banana" "Cherry"
// taken from https://github.com/charmbracelet/gum/tree/main/choose
package choose

import (
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
)

type model struct {
	Options
	Items        []item
	Selected     []item
	Quitting     bool
	Index        int
	numSelected  int
	currentOrder int
	paginator    paginator.Model
	aborted      bool
}

type item struct {
	text     string
	selected bool
	order    int
}

func (m model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		cmd = m.HandleKeys(msg)
		cmds = append(cmds, cmd)
	case ReturnSelectionsMsg:
		for _, item := range m.Items {
			if item.selected {
				m.Selected = append(m.Selected, item)
			}
		}
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	}

	m.paginator, cmd = m.paginator.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) HandleKeys(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	start, end := m.paginator.GetSliceBounds(len(m.Items))

	switch keypress := msg.String(); keypress {
	case "down", "j", "ctrl+j", "ctrl+n":
		m.Index++
		if m.Index >= len(m.Items) {
			m.Index = 0
			m.paginator.Page = 0
		}
		if m.Index >= end {
			m.paginator.NextPage()
		}
	case "up", "k", "ctrl+k", "ctrl+p":
		m.Index--
		if m.Index < 0 {
			m.Index = len(m.Items) - 1
			m.paginator.Page = m.paginator.TotalPages - 1
		}
		if m.Index < start {
			m.paginator.PrevPage()
		}
	case "right", "l", "ctrl+f":
		m.Index = clamp(m.Index+m.Height, 0, len(m.Items)-1)
		m.paginator.NextPage()
	case "left", "h", "ctrl+b":
		m.Index = clamp(m.Index-m.Height, 0, len(m.Items)-1)
		m.paginator.PrevPage()
	case "G":
		m.Index = len(m.Items) - 1
		m.paginator.Page = m.paginator.TotalPages - 1
	case "g":
		m.Index = 0
		m.paginator.Page = 0
	case "a":
		if m.Limit <= 1 {
			break
		}
		for i := range m.Items {
			if m.numSelected >= m.Limit {
				break // do not exceed given limit
			}
			if m.Items[i].selected {
				continue
			}
			m.Items[i].selected = true
			m.Items[i].order = m.currentOrder
			m.numSelected++
			m.currentOrder++
		}
	case "A":
		if m.Limit <= 1 {
			break
		}
		for i := range m.Items {
			m.Items[i].selected = false
			m.Items[i].order = 0
		}
		m.numSelected = 0
		m.currentOrder = 0
	case "ctrl+c", "esc", "q":
		m.aborted = true
		m.Quitting = true
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	case " ", "tab", "x":
		if m.Limit == 1 {
			break // no op
		}

		if m.Items[m.Index].selected {
			m.Items[m.Index].selected = false
			m.numSelected--
		} else if m.numSelected < m.Limit {
			m.Items[m.Index].selected = true
			m.Items[m.Index].order = m.currentOrder
			m.numSelected++
			m.currentOrder++
		}
	case "enter":
		m.Quitting = true
		// If the user hasn't selected any items in a multi-select.
		// Then we select the item that they have pressed enter on. If they
		// have selected items, then we simply return them.
		if m.numSelected < 1 {
			m.Items[m.Index].selected = true
		}
		cmd = ReturnSelectionsCmd()
		cmds = append(cmds, cmd)

	}

	return tea.Batch(cmds...)
}

type ReturnSelectionsMsg struct{}

func ReturnSelectionsCmd() tea.Cmd {
	return func() tea.Msg {
		return ReturnSelectionsMsg{}
	}
}

func (m model) View() string {
	//if m.quitting {
	//  return ""
	//}

	var s strings.Builder

	start, end := m.paginator.GetSliceBounds(len(m.Items))
	for i, item := range m.Items[start:end] {
		if i == m.Index%m.Height {
			s.WriteString(m.CursorStyle.Render(m.Cursor))
		} else {
			s.WriteString(strings.Repeat(" ", runewidth.StringWidth(m.Cursor)))
		}

		if item.selected {
			s.WriteString(m.SelectedItemStyle.Render(m.SelectedPrefix + item.text))
		} else if i == m.Index%m.Height {
			s.WriteString(m.CursorStyle.Render(m.CursorPrefix + item.text))
		} else {
			s.WriteString(m.ItemStyle.Render(m.UnselectedPrefix + item.text))
		}
		if i != m.Height {
			s.WriteRune('\n')
		}
	}

	if m.paginator.TotalPages <= 1 {
		return s.String()
	}

	s.WriteString(strings.Repeat("\n", m.Height-m.paginator.ItemsOnPage(len(m.Items))+1))
	s.WriteString("  " + m.paginator.View())

	return s.String()
}

//nolint:unparam
func clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
