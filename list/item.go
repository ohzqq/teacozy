package list

import (
	"strings"

	"github.com/ohzqq/teacozy/style"
	"github.com/sahilm/fuzzy"
)

type Item struct {
	fuzzy.Match
	Style            style.ListItem
	cursorPrefix     string
	selectedPrefix   string
	unselectedPrefix string
	isCur            bool
}

func NewItem(t string, idx int) Item {
	return Item{
		Match: fuzzy.Match{
			Str:   t,
			Index: idx,
		},
		Style: DefaultItemStyle(),
	}
}

func (match *Item) Render() string {
	var s strings.Builder
	mi := 0
	var buf strings.Builder
	for ci, c := range match.Str {
		// Check if the current character index matches the current matched index. If so, color the character to indicate a match.
		if mi < len(match.MatchedIndexes) && ci == match.MatchedIndexes[mi] {
			// Flush text buffer.
			s.WriteString(match.Style.Text.Render(buf.String()))
			buf.Reset()

			s.WriteString(match.Style.Match.Render(string(c)))
			// We have matched this character, so we never have to check it again. Move on to the next match.
			mi++
		} else {
			// Not a match, buffer a regular character.
			buf.WriteRune(c)
		}
	}
	// Flush text buffer.
	s.WriteString(match.Style.Text.Render(buf.String()))

	return s.String()
}

func (match *Item) RenderText() string {
	var s strings.Builder
	mi := 0
	var buf strings.Builder
	for ci, c := range match.Str {
		// Check if the current character index matches the current matched index. If so, color the character to indicate a match.
		if mi < len(match.MatchedIndexes) && ci == match.MatchedIndexes[mi] {
			// Flush text buffer.
			s.WriteString(match.Style.Text.Render(buf.String()))
			buf.Reset()

			s.WriteString(match.Style.Match.Render(string(c)))
			// We have matched this character, so we never have to check it again. Move on to the next match.
			mi++
		} else {
			// Not a match, buffer a regular character.
			buf.WriteRune(c)
		}
	}
	// Flush text buffer.
	s.WriteString(match.Style.Text.Render(buf.String()))

	return s.String()
}

func ChoicesToMatch(options []string) []Item {
	matches := make([]Item, len(options))
	for i, option := range options {
		matches[i] = NewItem(option, i)
	}
	return matches
}

func exactMatches(search string, choices []Item) []Item {
	matches := []Item{}
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
