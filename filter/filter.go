// Package filter provides a fuzzy searching text input to allow filtering a
// list of options to select one option.
//
// By default it will list all the files (recursively) in the current directory
// for the user to choose one, but the script (or user) can provide different
// new-line separated options to choose from.
//
// I.e. let's pick from a list of gum flavors:
//
// $ cat flavors.text | gum filter
package filter

import (
	"fmt"
	"strings"

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

type Model struct {
	textinput    textinput.Model
	viewport     *viewport.Model
	matches      []fuzzy.Match
	cursor       int
	Items        []Item
	FilterKeys   func(m *Model) keymap.KeyMap
	ListKeys     func(m *Model) keymap.KeyMap
	selected     map[int]struct{}
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
	case tea.KeyMsg:
		//switch keypress := msg.String(); keypress {
		//case "ctrl+c", "esc":
		//  switch m.filterState {
		//  case Filtering:
		//    m.filterState = Unfiltered
		//    m.textinput.Blur()
		//  default:
		//    m.aborted = true
		//    m.quitting = true
		//    cmds = append(cmds, tea.Quit)
		//  }
		//}

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
	case ReturnSelectionsMsg:
		m.quitting = true
		fmt.Printf("return sels sel %+V\n", msg)
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	}

	// It's possible that filtering items have caused fewer matches. So, ensure that the selected index is within the bounds of the number of matches.
	m.cursor = clamp(0, len(m.matches)-1, m.cursor)
	return m, tea.Batch(cmds...)
}

func ListKeyMap(m *Model) keymap.KeyMap {
	return keymap.KeyMap{
		keymap.NewBinding(
			keymap.WithKeys(" "),
			keymap.WithHelp("space", "select item"),
			keymap.WithCmd(SelectItemCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("down", "j"),
			keymap.WithHelp("down/j", "move cursor down"),
			keymap.WithCmd(DownCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("up", "k"),
			keymap.WithHelp("up/k", "move cursor up"),
			keymap.WithCmd(UpCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("ctrl+c", "esc", "q"),
			keymap.WithHelp("ctrl+c/esc/q", "quit"),
			keymap.WithCmd(tea.Quit),
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

func FilterKeyMap(m *Model) keymap.KeyMap {
	//start, end := m.paginator.GetSliceBounds(len(m.Items))
	return keymap.KeyMap{
		keymap.NewBinding(
			keymap.WithKeys("down", "ctrl+j"),
			keymap.WithHelp("down/ctrl+j", "move cursor down"),
			keymap.WithCmd(DownCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("up", "ctrl+k"),
			keymap.WithHelp("up/ctrl+k", "move cursor up"),
			keymap.WithCmd(UpCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("tab"),
			keymap.WithHelp("tab", "select item"),
			keymap.WithCmd(SelectItemCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("esc"),
			keymap.WithHelp("esc", "stop filtering"),
			keymap.WithCmd(StopFilteringCmd(m)),
		),
		keymap.NewBinding(
			keymap.WithKeys("ctrl+c"),
			keymap.WithHelp("ctrl+c", "quit"),
			keymap.WithCmd(tea.Quit),
		),
		keymap.NewBinding(
			keymap.WithKeys("enter"),
			keymap.WithHelp("enter", "return selections"),
			keymap.WithCmd(EnterCmd(m)),
		),
	}
}

func EnterCmd(m *Model) tea.Cmd {
	fmt.Println("enter")
	return ReturnSelectionsCmd(m)
}

func FilterItemsCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.filterState = Filtering
		m.cursor = 0
		m.textinput.Focus()
		return textinput.Blink()
	}
}

func StopFilteringCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.filterState = Unfiltered
		m.textinput.Blur()
		return nil
	}
}

type ReturnSelectionsMsg struct {
	choices []string
}

func ReturnSelectionsCmd(m *Model) tea.Cmd {
	fmt.Println("sels")
	return func() tea.Msg {
		if m.numSelected < 1 {
			m.Items[m.cursor].Selected = true
		}
		var sel ReturnSelectionsMsg
		for _, item := range m.Items {
			if item.Selected {
				sel.choices = append(sel.choices, m.Choices[item.Index])
			}
		}
		return sel
	}
}

func SelectItemCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.Limit == 1 {
			return nil
		}

		if _, ok := m.selected[m.matches[m.cursor].Index]; ok {
			delete(m.selected, m.matches[m.cursor].Index)
			m.numSelected--
		} else if m.numSelected < m.Limit {
			m.currentOrder++
			m.selected[m.matches[m.cursor].Index] = struct{}{}
			m.numSelected++
		}

		return nil
	}
}

func UpCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.Reverse {
			m.cursor = clamp(0, len(m.matches)-1, m.cursor+1)
			if len(m.matches)-m.cursor <= m.viewport.YOffset {
				m.viewport.SetYOffset(len(m.matches) - m.cursor - 1)
			}
		} else {
			m.cursor = clamp(0, len(m.matches)-1, m.cursor-1)
			if m.cursor < m.viewport.YOffset {
				m.viewport.SetYOffset(m.cursor)
			}
		}
		return nil
	}
}

func DownCmd(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.Reverse {
			m.cursor = clamp(0, len(m.matches)-1, m.cursor-1)
			if len(m.matches)-m.cursor > m.viewport.Height+m.viewport.YOffset {
				m.viewport.LineDown(1)
			}
		} else {
			m.cursor = clamp(0, len(m.matches)-1, m.cursor+1)
			if m.cursor >= m.viewport.YOffset+m.viewport.Height {
				m.viewport.LineDown(1)
			}
		}
		return nil
	}
}

func (m *Model) handleFilter(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)

	// yOffsetFromBottom is the number of lines from the bottom of the list to the top of the viewport. This is used to keep the viewport at a constant position when the number of matches are reduced in the reverse layout.
	var yOffsetFromBottom int
	if m.Reverse {
		yOffsetFromBottom = max(0, len(m.matches)-m.viewport.YOffset)
	}

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

	// For reverse layout, we need to offset the viewport so that the it remains at a constant position relative to the cursor.
	if m.Reverse {
		maxYOffset := max(0, len(m.matches)-m.viewport.Height)
		m.viewport.YOffset = clamp(0, maxYOffset, len(m.matches)-yOffsetFromBottom)
	}
	return cmd
}

func (m *Model) CursorUp() {
	if m.Reverse {
		m.cursor = clamp(0, len(m.matches)-1, m.cursor+1)
		if len(m.matches)-m.cursor <= m.viewport.YOffset {
			m.viewport.SetYOffset(len(m.matches) - m.cursor - 1)
		}
	} else {
		m.cursor = clamp(0, len(m.matches)-1, m.cursor-1)
		if m.cursor < m.viewport.YOffset {
			m.viewport.SetYOffset(m.cursor)
		}
	}
}

func (m *Model) CursorDown() {
	if m.Reverse {
		m.cursor = clamp(0, len(m.matches)-1, m.cursor-1)
		if len(m.matches)-m.cursor > m.viewport.Height+m.viewport.YOffset {
			m.viewport.LineDown(1)
		}
	} else {
		m.cursor = clamp(0, len(m.matches)-1, m.cursor+1)
		if m.cursor >= m.viewport.YOffset+m.viewport.Height {
			m.viewport.LineDown(1)
		}
	}
}

//func (m *Model) ToggleSelection() {
//  if _, ok := m.selected[m.matches[m.cursor].Str]; ok {
//    delete(m.selected, m.matches[m.cursor].Str)
//    m.numSelected--
//  } else if m.numSelected < m.Limit {
//    m.selected[m.matches[m.cursor].Str] = struct{}{}
//    m.numSelected++
//  }
//}

func (m Model) Init() tea.Cmd { return nil }
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var s strings.Builder

	// For reverse layout, if the number of matches is less than the viewport
	// height, we need to offset the matches so that the first match is at the
	// bottom edge of the viewport instead of in the middle.
	if m.Reverse && len(m.matches) < m.viewport.Height {
		s.WriteString(strings.Repeat("\n", m.viewport.Height-len(m.matches)))
	}

	// Since there are matches, display them so that the user can see, in real
	// time, what they are searching for.
	last := len(m.matches) - 1
	for i := range m.matches {
		// For reverse layout, the matches are displayed in reverse order.
		if m.Reverse {
			i = last - i
		}
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
			s.WriteString(m.UnselectedPrefixStyle.Render(m.SelectedPrefix))
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
	header := m.HeaderStyle.Render(m.Header)
	if m.Reverse {
		view := m.viewport.View() + "\n" + m.textinput.View()
		if m.Header != "" {
			return lipgloss.JoinVertical(lipgloss.Left, view, header)
		}

		return view
	}

	view := m.textinput.View() + "\n" + m.viewport.View()
	return lipgloss.JoinVertical(lipgloss.Left, header, view)
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
	for i, choice := range choices {
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
				Index:          i,
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
