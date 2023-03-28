package list

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/keys"
	"github.com/ohzqq/teacozy/style"
	"github.com/ohzqq/teacozy/util"
	"golang.org/x/exp/slices"
)

// FilterState describes the current filtering state on the model.
type FilterState int

// Possible filter states.
const (
	Unfiltered FilterState = iota // no filter set
	Filtering                     // user is actively setting a filter
)

type Model struct {
	Choices          []string
	choiceMap        []map[string]string
	Input            textinput.Model
	Viewport         *viewport.Model
	Paginator        paginator.Model
	Matches          []Item
	Items            []Item
	FilterKeys       func(m *Model) keys.KeyMap
	ListKeys         func(m *Model) keys.KeyMap
	Selected         map[int]struct{}
	numSelected      int
	cursor           int
	limit            int
	filterState      FilterState
	aborted          bool
	quitting         bool
	cursorPrefix     string
	selectedPrefix   string
	unselectedPrefix string
	header           string
	Placeholder      string
	Prompt           string
	Width            int
	Height           int
	Style            style.List
}

func New(choices ...string) *Model {
	tm := Model{
		Choices:          choices,
		Selected:         make(map[int]struct{}),
		FilterKeys:       FilterKeyMap,
		ListKeys:         ListKeyMap,
		filterState:      Unfiltered,
		Style:            DefaultStyle(),
		limit:            1,
		Prompt:           style.PromptPrefix,
		cursorPrefix:     style.CursorPrefix,
		selectedPrefix:   style.SelectedPrefix,
		unselectedPrefix: style.UnselectedPrefix,
		Height:           10,
	}

	w, h := util.TermSize()
	if tm.Height == 0 {
		tm.Height = h - 4
	}
	if tm.Width == 0 {
		tm.Width = w
	}
	tm.Input = textinput.New()
	tm.Input.Prompt = tm.Prompt
	tm.Input.PromptStyle = tm.Style.Prompt
	tm.Input.Placeholder = tm.Placeholder

	tm.header = "poot"

	return &tm
}

func (m *Model) Run() []int {
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	return m.Chosen()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.Height == 0 || m.Height > msg.Height {
			m.Viewport.Height = msg.Height - lipgloss.Height(m.Input.View())
		}

		// Make place in the view port if header is set
		if m.header != "" {
			m.Viewport.Height = m.Viewport.Height - lipgloss.Height(m.Style.Header.Render(m.header))
		}
		m.Viewport.Width = msg.Width
	case ReturnSelectionsMsg:
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		for _, k := range GlobalKeyMap(m) {
			if k.Matches(msg) {
				cmd = k.Cmd
				cmds = append(cmds, cmd)
			}
		}
		switch m.filterState {
		case Unfiltered:
			for _, k := range m.ListKeys(m) {
				if k.Matches(msg) {
					cmd = k.Cmd
					cmds = append(cmds, cmd)
				}
			}
		case Filtering:
			for _, k := range m.FilterKeys(m) {
				if k.Matches(msg) {
					cmd = k.Cmd
					cmds = append(cmds, cmd)
				}
			}
		}
		m.Input, cmd = m.Input.Update(msg)
		m.Matches = exactMatches(m.Input.Value(), m.Items)
		// If the search field is empty, let's not display the matches (none), but rather display all possible choices.
		if m.Input.Value() == "" {
			m.Matches = m.Items
		}
		cmds = append(cmds, cmd)
	}

	// It's possible that filtering items have caused fewer matches. So, ensure that the selected index is within the bounds of the number of matches.
	switch m.filterState {
	case Filtering:
		m.cursor = clamp(0, len(m.Matches)-1, m.cursor)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) CursorUp() {
	start, _ := m.Paginator.GetSliceBounds(len(m.Items))
	switch m.filterState {
	case Unfiltered:
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(m.Items) - 1
			m.Paginator.Page = m.Paginator.TotalPages - 1
		}
		if m.cursor < start {
			m.Paginator.PrevPage()
		}
	case Filtering:
		m.cursor = clamp(0, len(m.Matches)-1, m.cursor-1)
		if m.cursor < m.Viewport.YOffset {
			m.Viewport.SetYOffset(m.cursor)
		}
	}
}

func (m *Model) CursorDown() {
	_, end := m.Paginator.GetSliceBounds(len(m.Items))
	switch m.filterState {
	case Unfiltered:
		m.cursor++
		if m.cursor >= len(m.Items) {
			m.cursor = 0
			m.Paginator.Page = 0
		}
		if m.cursor >= end {
			m.Paginator.NextPage()
		}
	case Filtering:
		m.cursor = clamp(0, len(m.Matches)-1, m.cursor+1)
		if m.cursor >= m.Viewport.YOffset+m.Viewport.Height {
			m.Viewport.LineDown(1)
		}
	}
}

func (m *Model) ToggleSelection() {
	var idx int
	switch m.filterState {
	case Unfiltered:
		idx = m.Items[m.cursor].Index
	case Filtering:
		idx = m.Matches[m.cursor].Index
	}
	if _, ok := m.Selected[idx]; ok {
		delete(m.Selected, idx)
		m.numSelected--
		m.CursorDown()
		m.Items[idx].Toggle()
	} else if m.numSelected < m.limit {
		m.Selected[idx] = struct{}{}
		m.numSelected++
		m.CursorDown()
		m.Items[idx].Toggle()
	}
}

func (m Model) ItemIndex(c string) int {
	return slices.Index(m.Choices, c)
}

func (m Model) View() string {
	switch m.filterState {
	case Filtering:
		return m.FilteringView()
	default:
		return m.UnfilteredView()
	}
}

func (m Model) UnfilteredView() string {
	var s strings.Builder

	start, end := m.Paginator.GetSliceBounds(len(m.Items))

	for i, match := range m.Items[start:end] {
		if i == m.cursor%m.Height {
			match.IsCur()
		} else {
			match.NotCur()
		}

		s.WriteString(match.RenderPrefix())
		s.WriteString(match.Str)
		s.WriteRune('\n')
	}

	var view string
	if m.Paginator.TotalPages <= 1 {
		view = s.String()
	} else if m.Paginator.TotalPages > 1 {
		s.WriteString(strings.Repeat("\n", m.Height-m.Paginator.ItemsOnPage(len(m.Items))+1))
		s.WriteString("  " + m.Paginator.View())
	}

	view = s.String()
	if m.header != "" {
		header := m.Style.Header.Render(m.header + strings.Repeat(" ", m.Width))
		view = lipgloss.JoinVertical(lipgloss.Left, header, view)
	}

	m.Viewport.SetContent(view)
	view = m.Viewport.View()

	return view
}

func (m Model) FilteringView() string {
	var s strings.Builder

	for i, match := range m.Matches {
		switch {
		case i == m.cursor:
			match.IsCur()
		default:
			match.NotCur()
		}

		s.WriteString(match.RenderPrefix())
		text := lipgloss.StyleRunes(match.Str, match.MatchedIndexes, m.Style.Match, m.Style.Text)
		s.WriteString(text)

		s.WriteRune('\n')
	}

	m.Viewport.SetContent(s.String())

	view := m.Input.View() + "\n" + m.Viewport.View()
	return view
}

//nolint:unparam
func clamp(min, max, val int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (tm *Model) Init() tea.Cmd {
	tm.Items = ChoicesToMatch(tm.Choices)
	tm.Matches = tm.Items

	tm.Input.Width = tm.Width

	v := viewport.New(tm.Width, tm.Height+4)
	tm.Viewport = &v

	tm.Paginator = paginator.New()
	tm.Paginator.Type = paginator.Dots
	tm.Paginator.ActiveDot = style.Subdued.Render(style.Bullet)
	tm.Paginator.InactiveDot = style.VerySubdued.Render(style.Bullet)
	tm.Paginator.SetTotalPages((len(tm.Items) + tm.Height - 1) / tm.Height)
	tm.Paginator.PerPage = tm.Height
	return nil
}
