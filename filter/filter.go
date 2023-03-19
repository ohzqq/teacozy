package filter

import (
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ohzqq/teacozy/keymap"
	"github.com/ohzqq/teacozy/style"
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
	textinput    textinput.Model
	viewport     *viewport.Model
	paginator    paginator.Model
	matches      fuzzy.Matches
	cursor       int
	Items        fuzzy.Matches
	FilterKeys   func(m *Model) keymap.KeyMap
	ListKeys     func(m *Model) keymap.KeyMap
	selected     map[int]struct{}
	Choices      []string
	Chosen       []string
	numSelected  int
	filterState  FilterState
	currentOrder int
	aborted      bool
	quitting     bool
	Options
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.Height == 0 || m.Height > msg.Height {
			m.viewport.Height = msg.Height - lipgloss.Height(m.textinput.View())
		}

		// Make place in the view port if header is set
		if m.Header != "" {
			m.viewport.Height = m.viewport.Height - lipgloss.Height(m.HeaderStyle.Render(m.Header))
		}
		m.viewport.Width = msg.Width
	case ReturnSelectionsMsg:
		m.Chosen = msg.choices
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
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

	m.matches = exactMatches(m.textinput.Value(), m.Items)

	// If the search field is empty, let's not display the matches (none), but rather display all possible choices.
	if m.textinput.Value() == "" {
		m.matches = m.Items
	}

	return cmd
}

func (m *Model) CursorUp() {
	start, _ := m.paginator.GetSliceBounds(len(m.Items))
	switch m.filterState {
	case Unfiltered:
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(m.Items) - 1
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
	_, end := m.paginator.GetSliceBounds(len(m.Items))
	switch m.filterState {
	case Unfiltered:
		m.cursor++
		if m.cursor >= len(m.Items) {
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
		idx = m.Items[m.cursor].Index
	case Filtering:
		idx = m.matches[m.cursor].Index
	}
	if _, ok := m.selected[idx]; ok {
		delete(m.selected, idx)
		m.numSelected--
		m.CursorDown()
	} else if m.numSelected < m.Limit {
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

	start, end := m.paginator.GetSliceBounds(len(m.Items))
	items := m.renderItems(m.Items[start:end])
	s.WriteString(items)

	var view string
	if m.paginator.TotalPages <= 1 {
		view = s.String()
	}

	s.WriteString(strings.Repeat("\n", m.Height-m.paginator.ItemsOnPage(len(m.Items))+1))
	s.WriteString("  " + m.paginator.View())

	view = s.String()

	if m.Header != "" {
		header := m.HeaderStyle.Render(m.Header + strings.Repeat(" ", m.Width))
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
	if m.Limit == 1 {
		curPre = style.PromptPrefix
	}
	for i, match := range matches {
		var isCur bool
		switch {
		case m.filterState == Unfiltered && i == m.cursor%m.Height:
			fallthrough
		case m.filterState == Filtering && i == m.cursor:
			isCur = true
		}

		switch {
		case m.Limit > 1:
			s.WriteString("[")
		case m.Limit == 1:
			if !isCur {
				s.WriteString(strings.Repeat(" ", lipgloss.Width(curPre)))
			}
		}

		if isCur {
			s.WriteString(m.CursorStyle.Render(curPre))
		} else {
			if _, ok := m.selected[match.Index]; ok {
				s.WriteString(m.SelectedPrefixStyle.Render(m.SelectedPrefix))
			} else if m.Limit > 1 && !isCur {
				s.WriteString(m.UnselectedPrefixStyle.Render(m.UnselectedPrefix))
			}
		}
		if m.Limit > 1 {
			s.WriteString("]")
		}

		mi := 0
		var buf strings.Builder
		for ci, c := range match.Str {
			// Check if the current character index matches the current matched index. If so, color the character to indicate a match.
			if mi < len(match.MatchedIndexes) && ci == match.MatchedIndexes[mi] {
				// Flush text buffer.
				s.WriteString(m.TextStyle.Render(buf.String()))
				buf.Reset()

				s.WriteString(m.MatchStyle.Render(string(c)))
				// We have matched this character, so we never have to check it again. Move on to the next match.
				mi++
			} else {
				// Not a match, buffer a regular character.
				buf.WriteRune(c)
			}
		}
		// Flush text buffer.
		s.WriteString(m.TextStyle.Render(buf.String()))

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
