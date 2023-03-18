package choose

import (
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
	"github.com/ohzqq/teacozy/keymap"
	"github.com/sahilm/fuzzy"
)

type Model struct {
	Options
	textinput    textinput.Model
	viewport     *viewport.Model
	Items        []Item
	matches      []fuzzy.Match
	Selected     []map[int]string
	Quitting     bool
	cursor       int
	selected     map[string]struct{}
	KeyMap       func(m *Model) keymap.KeyMap
	Index        int
	numSelected  int
	showFilter   bool
	filterState  FilterState
	currentOrder int
	paginator    paginator.Model
	aborted      bool
}

type Item struct {
	Id       int
	Text     string
	Selected bool
	Order    int
}

func (m Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc", "q":
			m.aborted = true
			m.Quitting = true
			cmds = append(cmds, tea.Quit)
		}
		for _, k := range m.KeyMap(m) {
			if k.Matches(msg) {
				cmd := k.Cmd
				cmds = append(cmds, cmd)
			}
		}
	case ReturnSelectionsMsg:
		m.Quitting = true
		// If the user hasn't selected any items in a multi-select.
		// Then we select the item that they have pressed enter on. If they
		// have selected items, then we simply return them.
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	}

	m.paginator, cmd = m.paginator.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func FilterKeyMap(m *Model) keymap.KeyMap {
	start, end := m.paginator.GetSliceBounds(len(m.Items))
	return keymap.KeyMap{
		keymap.NewBinding(
			keymap.WithKeys("down", "ctrl+j"),
			keymap.WithHelp("down/ctrl+j", "move cursor down"),
			keymap.WithCmd(DownCmd(m, end)),
		),
		keymap.NewBinding(
			keymap.WithKeys("up", "ctrl+k"),
			keymap.WithHelp("up/ctrl+k", "move cursor up"),
			keymap.WithCmd(UpCmd(m, start)),
		),
		keymap.NewBinding(
			keymap.WithKeys("tab"),
			keymap.WithHelp("tab", "select item"),
			keymap.WithCmd(SelectItemCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("ctrl+c", "esc", "q"),
			keymap.WithHelp("ctrl+c/esc/q", "quit"),
		),
		keymap.NewBinding(
			keymap.WithKeys("enter"),
			keymap.WithHelp("enter", "return selections"),
			keymap.WithCmd(EnterCmd(m)),
		),
	}
}

func ListKeyMap(m *Model) keymap.KeyMap {
	start, end := m.paginator.GetSliceBounds(len(m.Items))
	return keymap.KeyMap{
		keymap.NewBinding(
			keymap.WithKeys("v"),
			keymap.WithHelp("v", "select all"),
			keymap.WithCmd(SelectAllItemsCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("V"),
			keymap.WithHelp("V", "deselect all"),
			keymap.WithCmd(DeselectAllItemsCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys(" "),
			keymap.WithHelp("space", "select item"),
			keymap.WithCmd(SelectItemCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("down", "j"),
			keymap.WithHelp("down/j", "move cursor down"),
			keymap.WithCmd(DownCmd(m, end)),
		),
		keymap.NewBinding(
			keymap.WithKeys("up", "k"),
			keymap.WithHelp("up/k", "move cursor up"),
			keymap.WithCmd(UpCmd(m, start)),
		),
		keymap.NewBinding(
			keymap.WithKeys("right", "l"),
			keymap.WithHelp("right/l", "next page"),
			keymap.WithCmd(NextPageCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("left", "h"),
			keymap.WithHelp("left/h", "prev page"),
			keymap.WithCmd(PrevPageCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("G"),
			keymap.WithHelp("G", "last item"),
			keymap.WithCmd(BottomCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("g"),
			keymap.WithHelp("g", "first item"),
			keymap.WithCmd(TopCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("ctrl+c", "esc", "q"),
			keymap.WithHelp("ctrl+c/esc/q", "quit"),
		),
		keymap.NewBinding(
			keymap.WithKeys("enter"),
			keymap.WithHelp("enter", "return selections"),
			keymap.WithCmd(EnterCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("/"),
			keymap.WithHelp("/", "filter items"),
			keymap.WithCmd(FilterItemsCmd(m)),
		),
	}
}

type ReturnSelectionsMsg struct{}

func SelectItemCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.Limit == 1 {
			return nil
		}

		if m.Items[m.Index].Selected {
			m.Items[m.Index].Selected = false
			m.numSelected--
		} else if m.numSelected < m.Limit {
			m.Items[m.Index].Selected = true
			m.Items[m.Index].Order = m.currentOrder
			m.numSelected++
			m.currentOrder++
		}
		return nil
	}
}

func SelectAllItemsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.Limit <= 1 {
			return nil
		}
		for i := range m.Items {
			if m.numSelected >= m.Limit {
				break // do not exceed given limit
			}
			if m.Items[i].Selected {
				continue
			}
			m.Items[i].Selected = true
			m.Items[i].Order = m.currentOrder
			m.numSelected++
			m.currentOrder++
		}
		return nil
	}
}

func EnterCmd(m *Model) tea.Cmd {
	return ReturnSelectionsCmd(m)
}

func ExitFilterCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.filterState = Unfiltered
		return nil
	}
}

func FilterItemsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		println("filter")
		m.filterState = Filtering
		m.paginator.Page = 0
		m.cursor = 0
		m.textinput.Focus()
		return textinput.Blink()
	}
}

func DeselectAllItemsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.Limit <= 1 {
			return nil
		}
		for i := range m.Items {
			m.Items[i].Selected = false
			m.Items[i].Order = 0
		}
		m.numSelected = 0
		m.currentOrder = 0
		return nil
	}
}

func ReturnSelectionsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.numSelected < 1 {
			m.Items[m.Index].Selected = true
		}
		for _, item := range m.Items {
			if item.Selected {
				sel := map[int]string{item.Id: item.Text}
				m.Selected = append(m.Selected, sel)
			}
		}
		return ReturnSelectionsMsg{}
	}
}

func UpCmd(m *Model, start int) tea.Cmd {
	return func() tea.Msg {
		m.Index--
		if m.Index < 0 {
			m.Index = len(m.Items) - 1
			m.paginator.Page = m.paginator.TotalPages - 1
		}
		if m.Index < start {
			m.paginator.PrevPage()
		}
		return nil
	}
}

func DownCmd(m *Model, end int) tea.Cmd {
	return func() tea.Msg {
		m.Index++
		if m.Index >= len(m.Items) {
			m.Index = 0
			m.paginator.Page = 0
		}
		if m.Index >= end {
			m.paginator.NextPage()
		}
		return nil
	}
}

func NextPageCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.Index = clamp(m.Index+m.Height, 0, len(m.Items)-1)
		m.paginator.NextPage()
		return nil
	}
}

func PrevPageCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.Index = clamp(m.Index-m.Height, 0, len(m.Items)-1)
		m.paginator.PrevPage()
		return nil
	}
}

func TopCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.Index = 0
		m.paginator.Page = 0
		return nil
	}
}

func BottomCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.Index = len(m.Items) - 1
		m.paginator.Page = m.paginator.TotalPages - 1
		return nil
	}
}

func (m Model) View() string {
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

		if item.Selected {
			s.WriteString(m.SelectedItemStyle.Render(m.SelectedPrefix + item.Text))
		} else if i == m.Index%m.Height {
			s.WriteString(m.CursorStyle.Render(m.CursorPrefix + item.Text))
		} else {
			s.WriteString(m.ItemStyle.Render(m.UnselectedPrefix + item.Text))
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
