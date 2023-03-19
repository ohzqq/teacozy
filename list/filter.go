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
	header           string
	Placeholder      string
	Prompt           string
	width            int
	height           int
	Style            style.List
}

func New(choices []string) *Model {
	tm := Model{
		Choices:          choices,
		selected:         make(map[int]struct{}),
		FilterKeys:       FilterKeyMap,
		ListKeys:         ListKeyMap,
		filterState:      Unfiltered,
		Style:            DefaultStyle(),
		limit:            1,
		Prompt:           style.PromptPrefix,
		cursorPrefix:     style.CursorPrefix,
		selectedPrefix:   style.SelectedPrefix,
		unselectedPrefix: style.UnselectedPrefix,
		height:           10,
	}

	w, h := util.TermSize()
	if tm.height == 0 {
		tm.height = h - 4
	}
	if tm.width == 0 {
		tm.width = w
	}
	tm.items = choicesToMatch(tm.Choices)
	tm.matches = tm.items

	return &tm
}

func (m *Model) Run() []int {
	m.textinput = textinput.New()
	m.textinput.Prompt = m.Prompt
	m.textinput.PromptStyle = m.Style.Prompt
	m.textinput.Placeholder = m.Placeholder
	m.textinput.Width = m.width

	v := viewport.New(m.width, m.height)
	m.viewport = &v

	m.paginator = paginator.New()
	m.paginator.SetTotalPages((len(m.items) + m.height - 1) / m.height)
	m.paginator.PerPage = m.height
	m.paginator.Type = paginator.Dots
	m.paginator.ActiveDot = style.Subdued.Render(style.Bullet)
	m.paginator.InactiveDot = style.VerySubdued.Render(style.Bullet)

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
		if m.height == 0 || m.height > msg.Height {
			m.viewport.Height = msg.Height - lipgloss.Height(m.textinput.View())
		}

		// Make place in the view port if header is set
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
		cmds = append(cmds, m.handleFilter(msg))
	}

	// It's possible that filtering items have caused fewer matches. So, ensure that the selected index is within the bounds of the number of matches.
	switch m.filterState {
	case Filtering:
		m.cursor = clamp(0, len(m.matches)-1, m.cursor)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) handleFilter(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)

	m.matches = exactMatches(m.textinput.Value(), m.items)

	// If the search field is empty, let's not display the matches (none), but rather display all possible choices.
	if m.textinput.Value() == "" {
		m.matches = m.items
	}

	return cmd
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
	}

	view = s.String()

	if m.header != "" {
		header := m.Style.Header.Render(m.header + strings.Repeat(" ", m.width))
		return lipgloss.JoinVertical(lipgloss.Left, header, view)
	}
	return view
}

func (m Model) FilteringView() string {
	var s strings.Builder

	items := m.renderItems(m.matches)
	s.WriteString(items)

	m.viewport.SetContent(s.String())

	view := m.textinput.View() + "\n" + m.viewport.View()
	return view
}

func (m Model) renderItems(matches fuzzy.Matches) string {
	var s strings.Builder
	curPre := style.CursorPrefix
	if m.limit == 1 {
		curPre = style.PromptPrefix
	}
	for i, match := range matches {
		var isCur bool
		switch {
		case m.filterState == Unfiltered && i == m.cursor%m.height:
			fallthrough
		case m.filterState == Filtering && i == m.cursor:
			isCur = true
		}

		switch {
		case m.limit > 1:
			s.WriteString("[")
		case m.limit == 1:
			if !isCur {
				s.WriteString(strings.Repeat(" ", lipgloss.Width(curPre)))
			}
		}

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

		mi := 0
		var buf strings.Builder
		for ci, c := range match.Str {
			// Check if the current character index matches the current matched index. If so, color the character to indicate a match.
			if mi < len(match.MatchedIndexes) && ci == match.MatchedIndexes[mi] {
				// Flush text buffer.
				s.WriteString(m.Style.Text.Render(buf.String()))
				buf.Reset()

				s.WriteString(m.Style.Match.Render(string(c)))
				// We have matched this character, so we never have to check it again. Move on to the next match.
				mi++
			} else {
				// Not a match, buffer a regular character.
				buf.WriteRune(c)
			}
		}
		// Flush text buffer.
		s.WriteString(m.Style.Text.Render(buf.String()))

		// We have finished displaying the match with all of it's matched characters highlighted and the rest filled in. Move on to the next match.
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

func exactMatches(search string, choices fuzzy.Matches) fuzzy.Matches {
	matches := fuzzy.Matches{}
	for _, choice := range choices {
		search = strings.ToLower(search)
		matchedString := strings.ToLower(choice.Str)

		index := strings.Index(matchedString, search)
		if index >= 0 {
			for s := range search {
				choice.MatchedIndexes = append(choice.MatchedIndexes, index+s)
			}
			matches = append(matches, choice)
		}
	}

	return matches
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
