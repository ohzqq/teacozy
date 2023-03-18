package filter

import (
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/ohzqq/teacozy/keymap"
	"github.com/sahilm/fuzzy"
)

// FilterState describes the current filtering state on the model.
type FilterState int

// Possible filter states.
const (
	Unfiltered    FilterState = iota // no filter set
	Filtering                        // user is actively setting a filter
	FilterApplied                    // a filter is applied and user is not editing filter
)

var (
	subduedStyle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#847A85", Dark: "#979797"})
	verySubduedStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"})
)

type Model struct {
	textinput    textinput.Model
	viewport     *viewport.Model
	paginator    paginator.Model
	matches      []fuzzy.Match
	cursor       int
	Items        []Item
	FilterKeys   func(m *Model) keymap.KeyMap
	ListKeys     func(m *Model) keymap.KeyMap
	selected     map[int]struct{}
	Chosen       []string
	numSelected  int
	filterState  FilterState
	currentOrder int
	aborted      bool
	quitting     bool
	Options
}

type Item struct {
	Index    int
	Text     string
	Selected bool
	Order    int
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	//fmt.Printf("cursor %v\n", m.cursor)
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
		if m.Reverse {
			m.viewport.YOffset = clamp(0, len(m.matches), len(m.matches)-m.viewport.Height)
		}
	case ReturnSelectionsMsg:
		m.Chosen = msg.choices
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch m.filterState {
		case Unfiltered:
			for _, k := range m.ListKeys(m) {
				if k.Matches(msg) {
					//fmt.Println(msg.String())
					cmd = k.Cmd
					cmds = append(cmds, cmd)
				}
			}
		case Filtering:
			for _, k := range m.FilterKeys(m) {
				if k.Matches(msg) {
					//fmt.Println(msg.String())
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

	// A character was entered, this likely means that the text input has changed. This suggests that the matches are outdated, so update them.
	if m.Fuzzy {
		m.matches = fuzzy.Find(m.textinput.Value(), m.Choices)
	} else {
		m.matches = exactMatches(m.textinput.Value(), m.Items)
	}

	// If the search field is empty, let's not display the matches (none), but rather display all possible choices.
	if m.textinput.Value() == "" {
		m.matches = matchAll(m.Items)
	}

	return cmd
}

func (m *Model) CursorUp() {
	//println(m.cursor)

	start, _ := m.paginator.GetSliceBounds(len(m.Items))
	//fmt.Println(start)
	//fmt.Println(end)

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
	switch m.filterState {
	case Unfiltered:
		if _, ok := m.selected[m.Items[m.cursor].Index]; ok {
			delete(m.selected, m.Items[m.cursor].Index)
			m.numSelected--
		} else if m.numSelected < m.Limit {
			m.selected[m.Items[m.cursor].Index] = struct{}{}
			m.numSelected++
			m.CursorDown()
		}
	case Filtering:
		if _, ok := m.selected[m.matches[m.cursor].Index]; ok {
			delete(m.selected, m.matches[m.cursor].Index)
			m.numSelected--
		} else if m.numSelected < m.Limit {
			m.currentOrder++
			m.selected[m.matches[m.cursor].Index] = struct{}{}
			m.numSelected++
			m.CursorDown()
		}
	}
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
	for i, item := range m.Items[start:end] {
		if i == m.cursor%m.Height {
			s.WriteString(m.CursorStyle.Render(m.CursorPrefix))
		} else {
			s.WriteString(strings.Repeat(" ", runewidth.StringWidth(m.CursorPrefix)))
		}

		if _, ok := m.selected[item.Index]; ok {
			s.WriteString(m.SelectedPrefixStyle.Render(m.SelectedPrefix) + item.Text)
		} else if m.Limit > 1 {
			s.WriteString(m.UnselectedPrefixStyle.Render(m.UnselectedPrefix) + item.Text)
		} else {
			s.WriteString(" " + item.Text)
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

	header := m.HeaderStyle.Render(m.Header)
	view := s.String()
	return lipgloss.JoinVertical(lipgloss.Left, header, view)
}

func (m Model) FilteringView() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	// Since there are matches, display them so that the user can see, in real
	// time, what they are searching for.
	for i := range m.matches {
		// For reverse layout, the matches are displayed in reverse order.
		match := m.matches[i]

		// If this is the current selected index, we add a small indicator to
		// represent it. Otherwise, simply pad the string.
		if i == m.cursor {
			s.WriteString(m.CursorStyle.Render(m.CursorPrefix))
		} else {
			s.WriteString(strings.Repeat(" ", runewidth.StringWidth(m.CursorPrefix)))
		}

		//If there are multiple selections mark them, otherwise leave an empty space
		if _, ok := m.selected[match.Index]; ok {
			s.WriteString(m.SelectedPrefixStyle.Render(m.SelectedPrefix))
		} else if m.Limit > 1 {
			s.WriteString(m.UnselectedPrefixStyle.Render(m.UnselectedPrefix))
		} else {
			s.WriteString(" ")
		}

		// For this match, there are a certain number of characters that have
		// caused the match. i.e. fuzzy matching.
		// We should indicate to the users which characters are being matched.
		mi := 0
		var buf strings.Builder
		for ci, c := range match.Str {
			// Check if the current character index matches the current matched
			// index. If so, color the character to indicate a match.
			if mi < len(match.MatchedIndexes) && ci == match.MatchedIndexes[mi] {
				// Flush text buffer.
				s.WriteString(m.TextStyle.Render(buf.String()))
				buf.Reset()

				s.WriteString(m.MatchStyle.Render(string(c)))
				// We have matched this character, so we never have to check it
				// again. Move on to the next match.
				mi++
			} else {
				// Not a match, buffer a regular character.
				buf.WriteRune(c)
			}
		}
		// Flush text buffer.
		s.WriteString(m.TextStyle.Render(buf.String()))

		// We have finished displaying the match with all of it's matched
		// characters highlighted and the rest filled in.
		// Move on to the next match.
		s.WriteRune('\n')
	}

	m.viewport.SetContent(s.String())

	// View the input and the filtered choices

	view := m.textinput.View() + "\n" + m.viewport.View()
	return view
}

func matchAll(options []Item) []fuzzy.Match {
	matches := make([]fuzzy.Match, len(options))
	for i, option := range options {
		matches[i] = fuzzy.Match{Str: option.Text, Index: i}
	}
	return matches
}

func exactMatches(search string, choices []Item) []fuzzy.Match {
	matches := fuzzy.Matches{}
	for _, choice := range choices {
		search = strings.ToLower(search)
		matchedString := strings.ToLower(choice.Text)

		index := strings.Index(matchedString, search)
		if index >= 0 {
			matchedIndexes := []int{}
			for s := range search {
				matchedIndexes = append(matchedIndexes, index+s)
			}
			matches = append(matches, fuzzy.Match{
				Str:            choice.Text,
				Index:          choice.Index,
				MatchedIndexes: matchedIndexes,
			})
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
