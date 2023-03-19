package list

import (
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/color"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
	"golang.org/x/exp/maps"
)

func New(choices []string) *Model {
	tm := Model{
		Choices:     choices,
		selected:    make(map[int]struct{}),
		FilterKeys:  FilterKeyMap,
		ListKeys:    ListKeyMap,
		filterState: Unfiltered,
	}
	tm.NoLimit = true

	tm.Prompt = style.PromptPrefix
	tm.PromptStyle = style.Prompt

	tm.CursorPrefix = style.CursorPrefix
	tm.SelectedPrefix = style.SelectedPrefix
	tm.UnselectedPrefix = style.UnselectedPrefix

	tm.CursorStyle = style.Cursor
	tm.SelectedPrefixStyle = style.Selected
	tm.UnselectedPrefixStyle = style.Unselected

	tm.TextStyle = style.Foreground
	tm.MatchStyle = lipgloss.NewStyle().Foreground(color.Cyan())
	//o.HeaderStyle = lipgloss.NewStyle().Foreground(color.Background).Background(color.Purple)
	tm.HeaderStyle = lipgloss.NewStyle().Foreground(color.Purple())
	tm.Height = 4

	tm.textinput = textinput.New()
	//model.textinput.Focus()

	tm.textinput.Prompt = tm.Prompt
	tm.textinput.PromptStyle = tm.PromptStyle
	tm.textinput.Placeholder = tm.Placeholder
	tm.textinput.Width = tm.Width

	w, h := util.TermSize()
	if tm.Height == 0 {
		tm.Height = h - 4
	}
	if tm.Width == 0 {
		tm.Width = w
	}

	v := viewport.New(tm.Width, tm.Height)
	tm.viewport = &v

	tm.Items = choicesToMatch(tm.Choices)
	tm.matches = tm.Items

	if tm.NoLimit {
		tm.Limit = len(tm.Choices)
	}

	pager := paginator.New()
	pager.SetTotalPages((len(tm.Items) + tm.Height - 1) / tm.Height)
	pager.PerPage = tm.Height
	pager.Type = paginator.Dots
	pager.ActiveDot = style.Subdued.Render("•")
	pager.InactiveDot = style.VerySubdued.Render("•")

	tm.paginator = pager
	return &tm
}

func EnterCmd(m *Model) tea.Cmd {
	return ReturnSelectionsCmd(m)
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

type ReturnSelectionsMsg struct {
	choices []string
}

func ReturnSelectionsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		var sel ReturnSelectionsMsg
		if len(m.selected) > 0 {
			for k := range m.selected {
				sel.choices = append(sel.choices, m.Choices[k])
			}
		} else if len(m.matches) > m.cursor && m.cursor >= 0 {
			sel.choices = append(sel.choices, m.matches[m.cursor].Str)
		}
		return sel
	}
}

func SelectItemCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.Limit == 1 {
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
		m.cursor = len(m.Items) - 1
		m.paginator.Page = m.paginator.TotalPages - 1
		return nil
	}
}

func NextPageCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.cursor = clamp(0, len(m.Items)-1, m.cursor+m.Height)
		m.paginator.NextPage()
		return nil
	}
}

func PrevPageCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.cursor = clamp(0, len(m.Items)-1, m.cursor-m.Height)
		m.paginator.PrevPage()
		return nil
	}
}

func SelectAllItemsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.Limit <= 1 {
			return nil
		}
		for i := range m.matches {
			if m.numSelected >= m.Limit {
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
		if m.Limit <= 1 {
			return nil
		}

		maps.Clear(m.selected)
		m.numSelected = 0

		return nil
	}
}
