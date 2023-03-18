package choose

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// FilterState describes the current filtering state on the model.
type FilterState int

// Possible filter states.
const (
	Unfiltered    FilterState = iota // no filter set
	Filtering                        // user is actively setting a filter
	FilterApplied                    // a filter is applied and user is not editing filter
)

// String returns a human-readable string of the current filter state.
func (f FilterState) String() string {
	return [...]string{
		"unfiltered",
		"filtering",
		"filter applied",
	}[f]
}

func (m Model) FilterView() string {
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

		// If there are multiple selections mark them, otherwise leave an empty space
		if _, ok := m.selected[match.Str]; ok {
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
