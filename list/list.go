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
	"github.com/ohzqq/teacozy/util"
	"github.com/sahilm/fuzzy"
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
	textinput        textinput.Model
	viewport         *viewport.Model
	paginator        paginator.Model
	matches          fuzzy.Matches
	items            fuzzy.Matches
	FilterKeys       func(m *Model) keys.KeyMap
	ListKeys         func(m *Model) keys.KeyMap
	selected         map[int]struct{}
	numSelected      int
	cursor           int
	limit            int
	filterState      FilterState
	aborted          bool
	quitting         bool
	cursorPrefix     string
	selectedPrefix   string
	unselectedPrefix string
	promptPrefix     string
	header           string
	Placeholder      string
	width            int
	height           int
	Style            Style
}

func New(choices []string, opts ...Option) *Model {
	tm := Model{
		Choices:          choices,
		selected:         make(map[int]struct{}),
		FilterKeys:       FilterKeyMap,
		ListKeys:         ListKeyMap,
		filterState:      Unfiltered,
		Style:            DefaultStyle(),
		limit:            1,
		promptPrefix:     PromptPrefix,
		cursorPrefix:     CursorPrefix,
		selectedPrefix:   SelectedPrefix,
		unselectedPrefix: UnselectedPrefix,
		height:           10,
	}

	w, h := util.TermSize()
	if tm.height == 0 {
		tm.height = h - 4
	}
	if tm.width == 0 {
		tm.width = w
	}

	for _, opt := range opts {
		opt(&tm)
	}

	tm.items = filterItems("", tm.Choices)
	tm.matches = tm.items

	tm.textinput = textinput.New()
	tm.textinput.Prompt = tm.promptPrefix
	tm.textinput.PromptStyle = tm.Style.Prompt
	tm.textinput.Placeholder = tm.Placeholder
	tm.textinput.Width = tm.width

	v := viewport.New(tm.width, tm.height)
	tm.viewport = &v

	tm.paginator = paginator.New()
	tm.paginator.SetTotalPages((len(tm.items) + tm.height - 1) / tm.height)
	tm.paginator.PerPage = tm.height
	tm.paginator.Type = paginator.Dots
	tm.paginator.ActiveDot = Subdued.Render(Bullet)
	tm.paginator.InactiveDot = VerySubdued.Render(Bullet)

	return &tm
}

func (tm *Model) Choose() []int {
	p := tea.NewProgram(tm)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	return tm.Chosen()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.height == 0 || m.height > msg.Height {
			m.viewport.Height = msg.Height - lipgloss.Height(m.textinput.View())
		}

		// Make place in the viewport if header is set
		if m.header != "" {
			m.viewport.Height = m.viewport.Height - lipgloss.Height(m.Style.Header.Render(m.header))
		}
		m.viewport.Width = msg.Width
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
		m.textinput, cmd = m.textinput.Update(msg)
		m.matches = filterItems(m.textinput.Value(), m.Choices)

		cmds = append(cmds, cmd)
	}

	// It's possible that filtering items have caused fewer matches. So, ensure that the selected index is within the bounds of the number of matches.
	switch m.filterState {
	case Filtering:
		m.cursor = clamp(0, len(m.matches)-1, m.cursor)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) CursorUp() {
	start, _ := m.paginator.GetSliceBounds(len(m.items))
	switch m.filterState {
	case Unfiltered:
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(m.items) - 1
			m.paginator.Page = m.paginator.TotalPages - 1
		}
		if m.cursor < start {
			m.paginator.PrevPage()
		}
	case Filtering:
		m.cursor = clamp(0, len(m.matches)-1, m.cursor-1)
		if m.cursor < m.viewport.YOffset {
			m.viewport.SetYOffset(m.cursor)
		}
	}
}

func (m *Model) CursorDown() {
	_, end := m.paginator.GetSliceBounds(len(m.items))
	switch m.filterState {
	case Unfiltered:
		m.cursor++
		if m.cursor >= len(m.items) {
			m.cursor = 0
			m.paginator.Page = 0
		}
		if m.cursor >= end {
			m.paginator.NextPage()
		}
	case Filtering:
		m.cursor = clamp(0, len(m.matches)-1, m.cursor+1)
		if m.cursor >= m.viewport.YOffset+m.viewport.Height {
			m.viewport.LineDown(1)
		}
	}
}

func (m *Model) ToggleSelection() {
	var idx int
	switch m.filterState {
	case Unfiltered:
		idx = m.items[m.cursor].Index
	case Filtering:
		idx = m.matches[m.cursor].Index
	}
	if _, ok := m.selected[idx]; ok {
		delete(m.selected, idx)
		m.numSelected--
		m.CursorDown()
	} else if m.numSelected < m.limit {
		m.selected[idx] = struct{}{}
		m.numSelected++
		m.CursorDown()
	}
}

func (m Model) Chosen() []int {
	var chosen []int
	if m.quitting {
		return chosen
	} else if len(m.selected) > 0 {
		for k := range m.selected {
			chosen = append(chosen, k)
		}
	} else if len(m.matches) > m.cursor && m.cursor >= 0 {
		chosen = append(chosen, m.cursor)
	}
	return chosen
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

	start, end := m.paginator.GetSliceBounds(len(m.items))
	items := m.renderItems(m.items[start:end])
	s.WriteString(items)

	var view string
	if m.paginator.TotalPages <= 1 {
		view = s.String()
	} else if m.paginator.TotalPages > 1 {
		s.WriteString(strings.Repeat("\n", m.height-m.paginator.ItemsOnPage(len(m.items))+1))
		s.WriteString("  " + m.paginator.View())
		view = s.String()
	}

	if m.header != "" {
		header := m.Style.Header.Render(m.header + strings.Repeat(" ", m.width))
		return lipgloss.JoinVertical(lipgloss.Left, header, view)
	}
	return view
}

func (m Model) FilteringView() string {
	m.viewport.SetContent(m.renderItems(m.matches))
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.textinput.View(),
		m.viewport.View(),
	)
}

func (m Model) renderItems(matches fuzzy.Matches) string {
	var s strings.Builder

	curPre := m.cursorPrefix
	if m.limit == 1 {
		curPre = m.promptPrefix
	}

	for i, match := range matches {
		var isCur bool

		// Determine if item is current
		switch {
		case m.filterState == Unfiltered && i == m.cursor%m.height:
			fallthrough
		case m.filterState == Filtering && i == m.cursor:
			isCur = true
		}

		// Write prefix
		switch {
		case m.limit > 1:
			s.WriteString("[")
		case m.limit == 1:
			if !isCur {
				s.WriteString(strings.Repeat(" ", lipgloss.Width(curPre)))
			}
		}

		// Style prefix
		if isCur {
			s.WriteString(m.Style.Cursor.Render(curPre))
		} else {
			if _, ok := m.selected[match.Index]; ok {
				s.WriteString(m.Style.SelectedPrefix.Render(m.selectedPrefix))
			} else if m.limit > 1 && !isCur {
				s.WriteString(m.Style.UnselectedPrefix.Render(m.unselectedPrefix))
			}
		}

		if m.limit > 1 {
			s.WriteString("]")
		}

		// Style item
		text := lipgloss.StyleRunes(
			match.Str,
			match.MatchedIndexes,
			m.Style.Match,
			m.Style.Text,
		)
		s.WriteString(text)

		s.WriteRune('\n')
	}

	return s.String()
}

func choicesToMatch(options []string) fuzzy.Matches {
	matches := make(fuzzy.Matches, len(options))
	for i, option := range options {
		matches[i] = fuzzy.Match{Str: option, Index: i}
	}
	return matches
}

func filterItems(search string, items []string) fuzzy.Matches {
	if search == "" {
		return choicesToMatch(items)
	}
	return fuzzy.Find(search, items)
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

func (m Model) Init() tea.Cmd { return nil }
